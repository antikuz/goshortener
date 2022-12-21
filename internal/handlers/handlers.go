package handlers

import (
	"math/rand"
	"time"

	"github.com/antikuz/goshortener/internal/db"
	"github.com/antikuz/goshortener/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"

type handler struct {
	storage *db.Storage
	logger  *zap.Logger
}

func NewHandler(storage *db.Storage, logger *zap.Logger) *handler {
	return &handler{
		storage: storage,
		logger:  logger,
	}
}

func (h *handler) Register(router *gin.Engine) {
	router.GET("/:shortUrl", h.shortURLRedirect)
	router.POST("/shorti", h.shortURLCreate)
}

func (h *handler) shortURLRedirect(c *gin.Context) {
	redirectURL, err := h.storage.GetURL(c.Params.ByName("shortUrl"))
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.Redirect(302, redirectURL.Target_url)
}

func (h *handler) shortURLCreate(c *gin.Context) {
	if shortingURL, ok := c.GetQuery("shortingURL"); ok {
		urlHash := generateURL()
		urlModel := models.ShortenURL{
			Key: urlHash,
			Target_url: shortingURL,
			Is_active: true,
		}

		if err := h.storage.AddShortURL(urlModel); err != nil {
			h.logger.Sugar().Errorf("Cant add new short url to database, due to err: %v", err)
			c.AbortWithStatus(500)
			return
		}

		c.JSON(200, gin.H{
			"message":   "short url created",
			"short_url": urlHash,
		})
	}
}

func generateURL() string {
	rand.Seed(time.Now().UnixMicro())
	result := ""
	for i := 5; i > 0; i-- {
		result += string(chars[rand.Intn(len(chars))])
	}

	return result
}