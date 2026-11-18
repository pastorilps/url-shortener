package delivery

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pastorilps/url-shortener/shortener/domain"
	"github.com/pastorilps/url-shortener/shortener/repository"
	"github.com/pastorilps/url-shortener/shortener/usecase"
)

type ShortenerDelivery struct {
	createUrl *usecase.CreateShortURL
	getUrl    *usecase.GetOriginalURL
}

func NewShortenerDelivery(
	createURL *usecase.CreateShortURL,
	getURL *usecase.GetOriginalURL,
) *ShortenerDelivery {
	return &ShortenerDelivery{
		createUrl: createURL,
		getUrl:    getURL,
	}
}

type CreateShortURlRequest struct {
	URL string `json:"url"`
}

func (d *ShortenerDelivery) Create(c *echo.Context) error {
	var request CreateShortURlRequest

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request",
		})
	}

	shortURL, err := d.createUrl.Execute(
		c.Request().Context(),
		request.URL,
	)
	if err != nil {
		c.Logger().Error("CreateShortURL failed", "error", err)

		if errors.Is(err, domain.ErrInvalidURL) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "could not create short URL",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"code":         shortURL.Code,
		"original_url": shortURL.OriginalURL,
		"short_url":    "http://localhost:8080/" + shortURL.Code,
		"created_ar":   shortURL.CreatedAt,
	})
}

func (d *ShortenerDelivery) Redirect(c *echo.Context) error {
	code := c.Param("code")

	shortURL, err := d.getUrl.Execute(
		c.Request().Context(),
		code,
	)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCode) ||
			errors.Is(err, domain.ErrExpiredURL) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}

		if errors.Is(err, repository.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "could not retrieve URL",
		})
	}

	return c.Redirect(http.StatusFound, shortURL.OriginalURL)
}

func (d *ShortenerDelivery) Get(c *echo.Context) error {
	code := c.Param("code")

	shortURL, err := d.getUrl.Execute(
		c.Request().Context(),
		code,
	)

	if err != nil {
		if errors.Is(err, domain.ErrInvalidCode) ||
			errors.Is(err, domain.ErrExpiredURL) ||
			errors.Is(err, repository.ErrNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "could not retrieve URL",
		})
	}

	return c.JSON(http.StatusOK, shortURL)
}

func (d *ShortenerDelivery) RegisterRoutes(e *echo.Echo) {
	e.POST("/api/v1/shorten", d.Create)
	e.GET("/api/v1/urls/:code", d.Get)

	e.GET("/:code", d.Redirect)
}
