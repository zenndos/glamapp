package config

import (
	"os"
	"strings"
	"sync/atomic"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURI string `mapstructure:"mongodb_uri"`
	Database    string `mapstructure:"mongo_database"`

	AppPort   string `mapstructure:"app_port"`
	AppHost   string `mapstructure:"app_host"`
	JWTSecret string `mapstructure:"jwt_secret"`
}

var confAT atomic.Value

func setDefaults() {
	viper.SetDefault("mongodb_uri", "mongodb://localhost:27017")
	viper.SetDefault("mongo_database", "glamapp")
	viper.SetDefault("app_port", "3001")
	viper.SetDefault("app_host", "localhost")
	viper.SetDefault("jwt_secret", "")
}

func ReadConfig() *Config {
	if conf, ok := confAT.Load().(*Config); ok && conf != nil {
		return conf
	}

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
	config.DatabaseURI = os.ExpandEnv(config.DatabaseURI)

	if !strings.HasPrefix(config.DatabaseURI, "mongodb://") {
		log.Fatal().Msgf("Invalid mongodb_uri: must start with mongodb://: [%s]", config.DatabaseURI)
	}

	if config.Database == "" {
		log.Fatal().Msg("mongo_database is not set")
	}

	if config.AppPort == "" {
		log.Fatal().Msg("app_port is not set")
	}

	if config.JWTSecret == "" {
		log.Fatal().Msg("jwt_secret is not set")
	}
	confAT.Store(config)

	return &config
}
