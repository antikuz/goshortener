package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/antikuz/goshortener/internal/db"
	"github.com/antikuz/goshortener/internal/handlers"
	"github.com/antikuz/goshortener/pkg/config"
	"github.com/antikuz/goshortener/pkg/logging"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx, _ := signal.NotifyContext(
		context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	cfg := config.LoadConfig()
	logger := logging.GetLogger(cfg.LogLevel)

	logger.Info("Create connection to database")
	database := db.NewStorage("./test.sqlite3", logger)
	defer database.Close()

	go func() {
		for range time.Tick(time.Hour) {
			expiredURLs, err := database.GetExpiredURL()
			if err != nil {
				logger.Sugar().Errorf("Cannot get expired urls, due to err: %v", err)
			}
			for _, url := range expiredURLs {
				err := database.DeleteURL(url.Key)
				if err != nil {
					logger.Sugar().Errorf("Cannot delete expired url with hash %s, due to err: %v", url.Key, err)
				}
			}
		}
	}()

	router := gin.Default()
	handler := handlers.NewHandler(database, logger)
	handler.Register(router)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		logger.Sugar().Infof("Start goshortener at %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Sugar().Fatalf("Can't start webserver due to err: %v", err)
		}
	}()

	<- ctx.Done()
	logger.Info("Shutting down server...")
}