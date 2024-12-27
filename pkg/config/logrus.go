package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewLogrus(viper *viper.Viper) *logrus.Logger {
	log := logrus.New()

	appEnv := viper.GetString("APP_ENV")
	if appEnv == "production" || appEnv == "prod" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return log
}
