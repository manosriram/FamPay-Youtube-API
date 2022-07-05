package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	YoutubeApiKeys []string `mapstructure:"YOUTUBE_DEVELOPER_KEY" validate:"required"`
	MongoUsername  string   `mapstructure:"MONGO_USERNAME" validate:"required"`
	MongoPassword  string   `mapstructure:"MONGO_PASSWORD" validate:"required"`
}

func Load(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	config := Config{}
	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
