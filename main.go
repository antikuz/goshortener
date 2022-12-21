package main

import (
	"github.com/antikuz/goshortener/internal/db"
	"github.com/antikuz/goshortener/pkg/logging"
	"github.com/antikuz/goshortener/internal/handlers"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)






func main() {
	logger := logging.GetLogger()
	logger.Sugar().Info("Create connection to database")
	database := db.NewStorage("./test.sqlite3", logger)
	defer database.Close()

	router := gin.Default()
	handler := handlers.NewHandler(database, logger)
	handler.Register(router)

	logger.Sugar().Fatalf("Can't start webserver due to err: %v", router.Run())
}


