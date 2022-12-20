package main

import (
	"math/rand"
	"time"

	"github.com/antikuz/goshortener/internal/db"
	"github.com/antikuz/goshortener/pkg/logging"
	_ "github.com/mattn/go-sqlite3"
)

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890"


func generateURL() string {
	rand.Seed(time.Now().UnixMicro())
	result := ""
	for i := 5; i > 0; i-- {
		result += string(chars[rand.Intn(len(chars))])
	}

	return result
}

func main() {
	logger := logging.GetLogger()
	logger.Sugar().Info("Create connection to database")
	db.NewStorage("./test.sqlite3", logger)
}
