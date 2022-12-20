package logging

import (
	"go.uber.org/zap"
)

var logLevel = "debug"
var logger *zap.Logger

func GetLogger() *zap.Logger {
	return logger
}

func init() {
	if logLevel == "debug" {
		logger = zap.Must(zap.NewDevelopment())
	} else {
		logger = zap.Must(zap.NewProduction())
	}
}