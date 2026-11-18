package repository

import (
	"context"
	"errors"

	"github.com/pastorilps/url-shortener/shortener/entity"
)

var (
	ErrCodeAlreadyExists = errors.New("short code already exists")
	ErrNotFound          = errors.New("short code not found")
)

type ShortenerRepository interface {
	Create(ctx context.Context, shortUrl *entity.ShortUrl) error

	FindByCode(
		ctx context.Context,
		code string,
	) (*entity.ShortUrl, error)

	Delete(ctx context.Context, code string) error
}
