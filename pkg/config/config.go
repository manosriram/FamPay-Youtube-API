package config

import "github.com/spf13/viper"

type Config struct {
	YoutubeDeveloperKey string `mapstructure:"YOUTUBE_DEVELOPER_KEY" validate:"required"`
}

func Load(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")

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
