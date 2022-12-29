package logging

import (
	"log"
	"sync"

	"go.uber.org/zap"
)

var logger *zap.Logger
var once sync.Once

func GetLogger(logLevel string) *zap.Logger {
	once.Do(func() {
		if logLevel == "debug" {
			logger = zap.Must(zap.NewDevelopment())
		} else {
			level, err := zap.ParseAtomicLevel(logLevel)
			if err != nil {
				log.Fatal(err)
			}
			config := zap.Config{
				Level: level,
				Development: false,
				Sampling: &zap.SamplingConfig{
					Initial:    100,
					Thereafter: 100,
				},
				Encoding:         "json",
				EncoderConfig:    zap.NewProductionEncoderConfig(),
				OutputPaths:      []string{"stderr"},
				ErrorOutputPaths: []string{"stderr"},
			}

			logger = zap.Must(config.Build())
		}
	})

	return logger
}