package usecase

import (
	"context"

	"github.com/pastorilps/url-shortener/shortener/domain"
	"github.com/pastorilps/url-shortener/shortener/entity"
	"github.com/pastorilps/url-shortener/shortener/repository"
)

type GetOriginalURL struct {
	repository repository.ShortenerRepository
}

func NewGetOriginalURL(
	repo repository.ShortenerRepository,
) *GetOriginalURL {
	return &GetOriginalURL{
		repository: repo,
	}
}

func (uc *GetOriginalURL) Execute(
	ctx context.Context,
	code string,
) (*entity.ShortUrl, error) {
	if err := domain.ValidadeCode(code); err != nil {
		return nil, err
	}

	shortUrl, err := uc.repository.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if domain.IsExpired(shortUrl.ExpiresAt) {
		return nil, domain.ErrExpiredURL
	}

	return shortUrl, nil
}
