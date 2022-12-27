package config

import (
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type config struct {
	Host     string
	Port     int
	LogLevel string
}

func LoadConfig() *config {
	var c config

	pflag.Int("port", 8080, "help message for flagname")
	pflag.String("host", "", "help message for flagname")
	pflag.String("loglevel", "debug", "help message for flagname")
	pflag.Parse()

	viper.AutomaticEnv()
	viper.BindPFlags(pflag.CommandLine)

	err := viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode config into struct, %v", err)
	}

	return &c
}
