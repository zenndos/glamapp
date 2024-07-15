package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURI string `json:"mongodb_uri"`
	Database    string `json:"mongo_database"`

	APP_PORT string `json:"app_port"`
}

func setDefaults() {
	viper.SetDefault("mongodb_uri", "mongodb://localhost:27017")
	viper.SetDefault("app_port", "3000")
}

func ReadConfig() *Config {
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatal().Err(err).Msg("Error reading config file")
		}
	}

	setDefaults()

	config := Config{}
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal().Err(err).Msg("error unmarshalling config")
	}

	return &config
}
