package cassandra

import (
	"context"
	"errors"

	"github.com/gocql/gocql"
	"github.com/pastorilps/url-shortener/shortener/entity"
	"github.com/pastorilps/url-shortener/shortener/repository"
)

type ShortenerRepository struct {
	session *gocql.Session
}

func NewShortenerRepository(
	session *gocql.Session,
) *ShortenerRepository {
	return &ShortenerRepository{
		session: session,
	}
}

func (r *ShortenerRepository) Create(
	ctx context.Context,
	shortURL *entity.ShortUrl,
) error {

	applied, err := r.session.Query(`
        INSERT INTO short_urls (
            code,
            original_url,
            created_at,
            expires_at
        ) VALUES (?, ?, ?, ?)
        IF NOT EXISTS
    `,
		shortURL.Code,
		shortURL.OriginalURL,
		shortURL.CreatedAt,
		shortURL.ExpiresAt,
	).WithContext(ctx).ScanCAS()

	if err != nil {
		return err
	}

	if !applied {
		return repository.ErrCodeAlreadyExists
	}

	return nil
}

func (r *ShortenerRepository) FindByCode(
	ctx context.Context,
	code string,
) (*entity.ShortUrl, error) {

	shortURL := &entity.ShortUrl{}

	err := r.session.Query(`
        SELECT code, original_url, created_at, expires_at
        FROM short_urls
        WHERE code = ?
    `, code).
		WithContext(ctx).
		Scan(
			&shortURL.Code,
			&shortURL.OriginalURL,
			&shortURL.CreatedAt,
			&shortURL.ExpiresAt,
		)

	if errors.Is(err, gocql.ErrNotFound) {
		return nil, repository.ErrNotFound
	}

	if err != nil {
		return nil, err
	}

	return shortURL, nil
}

func (r *ShortenerRepository) Delete(
	ctx context.Context,
	code string,
) error {

	return r.session.Query(`
        DELETE FROM short_urls WHERE code = ?
    `, code).
		WithContext(ctx).
		Exec()
}

var _ repository.ShortenerRepository = (*ShortenerRepository)(nil)
