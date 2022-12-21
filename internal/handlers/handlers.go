package handlers

import (
	"github.com/antikuz/goshortener/internal/db"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

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

func (h *handler) ShortURLRedirect(c *gin.Context) {
	redirectURL, err := h.storage.GetURL(c.Params.ByName("shortUrl"))
	if err != nil {
		c.AbortWithStatus(404)
		return
	}
	
	c.String(200, "hash %s url: %v", c.Params.ByName("shortUrl"), redirectURL)
}

func (h *handler) ShortURLCreate(c *gin.Context) {
	c.JSON(200, gin.H{
		"message":   "short url created successfully",
		"short_url": "sdhsdh",
	})
}
