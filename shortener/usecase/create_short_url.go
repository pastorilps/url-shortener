package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/pastorilps/url-shortener/shortener/domain"
	"github.com/pastorilps/url-shortener/shortener/entity"
	"github.com/pastorilps/url-shortener/shortener/repository"
	"github.com/redis/go-redis/v9"
)

type CreateShortURL struct {
	repository repository.ShortenerRepository
}

func NewCreateShortURL(
	repo repository.ShortenerRepository,
) *CreateShortURL {
	return &CreateShortURL{
		repository: repo,
	}
}

func (uc *CreateShortURL) Execute(
	ctx context.Context,
	originalURL string,
) (*entity.ShortUrl, error) {
	if err := domain.ValidateURL(originalURL); err != nil {
		return nil, err
	}

	ctx = context.Background()

	addrs := []string{
		"localhost:7000",
		"localhost:7001",
		"localhost:7002",
		"localhost:7003",
		"localhost:7004",
		"localhost:7005",
	}

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
		ClusterSlots: func(ctx context.Context) ([]redis.ClusterSlot, error) {
			return getClusterSlots(ctx, addrs)
		},
	})

	defer rdb.Close()

	key := "counter:url"
	startValue := 14000000

	err := rdb.SetNX(ctx, key, startValue, 0).Err()
	if err != nil {
		return nil, fmt.Errorf("initialize URL counter: %w", err)
	}

	newId, err := rdb.Incr(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("increment URL counter: %w", err)
	}

	for attempt := 0; attempt < 5; attempt++ {
		code, err := domain.GenerateCode(uint64(newId))
		if err != nil {
			return nil, err
		}

		expireDate := time.Now().UTC().AddDate(5, 0, 0)

		shortUrl := &entity.ShortUrl{
			Code:        code,
			OriginalURL: originalURL,
			CreatedAt:   time.Now().UTC(),
			ExpiresAt:   &expireDate,
		}

		err = uc.repository.Create(ctx, shortUrl)

		if err == repository.ErrCodeAlreadyExists {
			continue
		}

		if err != nil {
			return nil, err
		}

		return shortUrl, nil
	}

	return nil, repository.ErrCodeAlreadyExists
}

func getClusterSlots(ctx context.Context, addrs []string) ([]redis.ClusterSlot, error) {
	clients := make([]*redis.Client, len(addrs))
	defer func() {
		for _, client := range clients {
			if client != nil {
				_ = client.Close()
			}
		}
	}()

	nodeAddrs := make(map[string]string, len(addrs))
	for i, addr := range addrs {
		client := redis.NewClient(&redis.Options{Addr: addr})
		clients[i] = client

		nodeID, err := client.ClusterMyID(ctx).Result()
		if err != nil {
			return nil, fmt.Errorf("get Redis cluster node ID from %s: %w", addr, err)
		}
		nodeAddrs[nodeID] = addr
	}

	shards, err := clients[0].ClusterShards(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("get Redis cluster shards: %w", err)
	}

	var slots []redis.ClusterSlot
	for _, shard := range shards {
		nodes := make([]redis.ClusterNode, 0, len(shard.Nodes))
		for _, node := range shard.Nodes {
			addr, ok := nodeAddrs[node.ID]
			if !ok {
				return nil, fmt.Errorf("Redis cluster node %s has no published address", node.ID)
			}
			nodes = append(nodes, redis.ClusterNode{
				ID:   node.ID,
				Addr: addr,
			})
		}

		for _, slotRange := range shard.Slots {
			slots = append(slots, redis.ClusterSlot{
				Start: int(slotRange.Start),
				End:   int(slotRange.End),
				Nodes: nodes,
			})
		}
	}

	return slots, nil
}
