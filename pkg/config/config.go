package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	YoutubeApiKeys []string `mapstructure:"YOUTUBE_DEVELOPER_KEY" validate:"required"`
	MongoURI       string   `mapstructure:"MONGO_URI" validate:"required"`
}

// Loads config.yaml into config.Config
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
