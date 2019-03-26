package utils

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// logger
var Log = log.WithFields(log.Fields{
	"env": os.Getenv("GO_ENV")})

// Find and read the config file
func ReadConfig(path string) *viper.Viper {
	var v = viper.New()
	v.SetConfigType("yaml")
	//v.SetConfigName("log")
	v.AddConfigPath(path)

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	return v
}
