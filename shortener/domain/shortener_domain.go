package domain

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/sqids/sqids-go"
)

var (
	ErrInvalidURL  = errors.New("invalid URL")
	ErrInvalidCode = errors.New("invalid short code")
	ErrExpiredURL  = errors.New("URL has expired")
)

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func generateAlphabet(secret string) string {
	hash := sha256.Sum256([]byte(secret))

	seed := int64(binary.BigEndian.Uint64(hash[:8]))

	r := rand.New(rand.NewSource(seed))

	runes := []rune(alphabet)
	r.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})

	return string(runes)
}

func ValidateURL(rawUrl string) error {
	parsed, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		return ErrInvalidURL
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ErrInvalidURL
	}

	if parsed.Host == "" {
		return ErrInvalidURL
	}

	return nil
}

func ValidadeCode(code string) error {
	if len(code) < 4 || len(code) > 12 {
		return ErrInvalidCode
	}

	for _, char := range code {
		if !strings.ContainsRune(alphabet, char) {
			return ErrInvalidCode
		}
	}

	return nil
}

func GenerateCode(id uint64) (string, error) {
	secret := "secret"

	customAlphabet := generateAlphabet(secret)

	s, err := sqids.New(sqids.Options{
		Alphabet: customAlphabet,
	})

	if err != nil {
		panic(err)
	}

	hashId, _ := s.Encode([]uint64{id})

	return hashId, nil
}

func IsExpired(expiresAt *time.Time) bool {
	if expiresAt == nil {
		return false
	}

	return !expiresAt.IsZero() && time.Now().After(*expiresAt)
}
