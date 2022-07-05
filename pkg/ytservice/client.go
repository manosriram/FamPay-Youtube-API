package ytservice

import (
	"context"
	"errors"
	"fmt"

	config "github.com/manosriram/youtubeAPI-fampay/pkg/config"
	option "google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func CreateYoutubeClient(config *config.Config) (*youtube.Service, error) {
	keys := config.YoutubeApiKeys
	for _, key := range keys {
		youtubeClient, err := youtube.NewService(context.TODO(), option.WithAPIKey(key))
		if err != nil {
			continue
		}

		return youtubeClient, nil
	}

	fmt.Println(keys)
	return nil, errors.New("invalid API key")
}
