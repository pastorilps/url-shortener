package entity

import "time"

type ShortUrl struct {
	Code        string
	OriginalURL string
	CreatedAt   time.Time
	ExpiresAt   *time.Time
}
