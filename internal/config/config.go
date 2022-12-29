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
	DBPath   string
}

func LoadConfig() *config {
	var c config

	pflag.String("host", "", "IP to listen on, defaults to all IPs")
	pflag.Int("port", 8080, "port to listen")
	pflag.String("dbpath", "./db.sqlite3", "path to db file")
	pflag.String("loglevel", "debug", "logger level")
	pflag.CommandLine.SortFlags = false
	pflag.Parse()


	viper.AutomaticEnv()
	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		log.Fatalf("unable to BindPFlags, due to error:%v", err)
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatalf("unable to decode config into struct, %v", err)
	}

	return &c
}
