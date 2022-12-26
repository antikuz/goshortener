package main

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/antikuz/goshortener/internal/db"
	"github.com/antikuz/goshortener/internal/handlers"
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

	logger := logging.GetLogger()
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

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Sugar().Fatalf("Can't start webserver due to err: %v", err)
		}
	}()

	<- ctx.Done()
	logger.Info("Shutting down server...")
}