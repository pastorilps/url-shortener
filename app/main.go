package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v5"
	"github.com/pastorilps/url-shortener/infrastructure/cassandra"
	"github.com/pastorilps/url-shortener/middlewares"
	"github.com/pastorilps/url-shortener/shortener/usecase"

	shortenerDelivery "github.com/pastorilps/url-shortener/shortener/delivery"
)

func main() {
	session, err := cassandra.NewConnection(
		cassandra.Config{
			Hosts: []string{
				getEnv("CASSANDRA_HOST", "localhost"),
			},
			Keyspace: getEnv("CASSANDRA_KEYSPACE", "application"),
			Username: os.Getenv("CASSANDRA_USERNAME"),
			Password: os.Getenv("CASSANDRA_PASSWORD"),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	defer session.Close()

	shortnerRepository := cassandra.NewShortenerRepository(session)

	createURL := usecase.NewCreateShortURL(shortnerRepository)
	getURL := usecase.NewGetOriginalURL(shortnerRepository)

	delivery := shortenerDelivery.NewShortenerDelivery(
		createURL,
		getURL,
	)

	e := echo.New()

	e.Use(middlewares.CORS())
	e.Use(middlewares.RateLimit())

	delivery.RegisterRoutes(e)

	port := getEnv("PORT", "8080")

	log.Printf("server running on port %s", port)

	if err := e.Start(":" + port); err != nil {
		log.Fatal(err)
	}

}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
