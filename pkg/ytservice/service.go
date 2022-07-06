package ytservice

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/manosriram/youtubeAPI-fampay/data"
	"github.com/manosriram/youtubeAPI-fampay/db"
	config "github.com/manosriram/youtubeAPI-fampay/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	option "google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Service struct {
	Config *config.Config
}

var (
	maxResults      = flag.Int64("max-results", 6, "Max YouTube results")
	predefinedQuery = "football"
)

/*
   pings youtube for video with search query, uses available API keys.
   picks the first valid API key.
*/
func getMetadataFromYoutube(logger *zap.SugaredLogger, config *config.Config, query string) (*youtube.SearchListResponse, error) {
	keys := config.YoutubeApiKeys
	for _, key := range keys {
		youtubeClient, err := youtube.NewService(context.TODO(), option.WithAPIKey(key))
		if err != nil {
			continue
		}
		call := youtubeClient.Search.List([]string{"id,snippet"}).
			Q(query).
			MaxResults(*maxResults).
			Order("date").
			Type("video").
			PublishedAfter("2020-01-01T00:00:00Z")

		response, err := call.Do()
		if err != nil {
			continue
		}

		logger.Infow("got video metadata from youtube")
		return response, nil
	}
	return nil, errors.New("invalid api key(s)")
}

/*
   fetches videos metadata from youtube with the predefinedQuery in descending order of published date and stores it in DB.
   supports multiple API keys, picks up the first valid key.
   bulk inserts the video array into DB.
*/

func FetchVideosByQuery(logger *zap.SugaredLogger, config *config.Config, query string, mongoCollection *mongo.Collection) error {
	response, err := getMetadataFromYoutube(logger, config, query)
	if err != nil {
		logger.Errorw("error getting video metadata from youtube", "error", err)
		return err
	}

	videosList := []data.Video{}

	// Create a list of videos metadata to upsert
	for _, item := range response.Items {
		newVideo := data.Video{
			Title:        item.Snippet.Title,
			Description:  item.Snippet.Description,
			PublishedAt:  item.Snippet.PublishedAt,
			ThumbnailUrl: item.Snippet.Thumbnails.Default.Url,
			VideoETag:    item.Etag,
		}
		videosList = append(videosList, newVideo)
	}

	err = db.BulkUpsertVideos(context.TODO(), logger, videosList, mongoCollection)
	if err != nil {
		logger.Errorw("error bulk upserting videos", "error", err)
		return err
	}
	fmt.Println(videosList)
	return nil
}

/*
   service which fetches videos metadata from DB and returns the array.
   returns results matching the query if query is specified.
*/
func (svc Service) LoadStoredVideos(ctx context.Context, showVideoRequest data.ShowVideoRequest, logger *zap.SugaredLogger, mongoCollection *mongo.Collection) ([]*data.Video, error) {
	videos, err := db.GetVideosList(ctx, showVideoRequest, logger, mongoCollection)
	if err != nil {
		return nil, err
	}
	return videos, nil
}
