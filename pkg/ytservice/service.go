package ytservice

import (
	"context"
	"flag"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/manosriram/youtubeAPI-fampay/data"
	"github.com/manosriram/youtubeAPI-fampay/db"
	config "github.com/manosriram/youtubeAPI-fampay/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type Service struct {
	Config *config.Config
}

var (
	maxResults      = flag.Int64("max-results", 6, "Max YouTube results")
	predefinedQuery = "football"
)

func printIDs(sectionName string, matches map[string]string) {
	fmt.Printf("%v:\n", sectionName)
	fmt.Println(matches)
	for id, title := range matches {
		fmt.Printf("[%v] %v\n", id, title)
	}
	fmt.Printf("\n\n")
}

// Fetches video(s) metadata from YouTube and stores in DB
func FetchVideosByQuery(logger *zap.SugaredLogger, config config.Config, query string, mongoCollection *mongo.Collection) error {
	client := &http.Client{
		Transport: &transport.APIKey{Key: config.YoutubeDeveloperKey},
	}

	youtubeClient, err := youtube.New(client)
	if err != nil {
		logger.Errorw("error creating youtube client", "error", err)
		return err
	}
	call := youtubeClient.Search.List([]string{"id,snippet"}).
		Q(query).
		MaxResults(*maxResults)
	response, err := call.Do()
	if err != nil {
		logger.Errorw("error calling youtube search API", "error", err)
		return err
	}

	videosList := []data.Video{}

	for _, item := range response.Items {
		switch item.Id.Kind {
		case "youtube#video":
			newVideoId := uuid.New()
			newVideo := data.Video{
				Id:           newVideoId,
				Title:        item.Snippet.Title,
				Description:  item.Snippet.Description,
				PublishedAt:  item.Snippet.PublishedAt,
				ThumbnailUrl: item.Snippet.Thumbnails.Default.Url,
				VideoETag:    item.Etag,
			}
			videosList = append(videosList, newVideo)
		}
	}
	err = db.BulkUpsertVideos(context.TODO(), logger, videosList, mongoCollection)
	if err != nil {
		logger.Errorw("error bulk upserting videos", "error", err)
		return err
	}
	fmt.Println(videosList)
	return nil
}

func (svc Service) LoadStoredVideos(ctx context.Context, logger *zap.SugaredLogger, mongoCollection *mongo.Collection) ([]*data.Video, error) {
	videos, err := db.GetVideosList(ctx, logger, mongoCollection)
	if err != nil {
		return nil, err
	}
	// FetchVideosByQuery(logger, predefinedQuery, mongoCollection)
	return videos, nil
}

func (svc Service) LoadStoredVideosByQuery(ctx context.Context) error {
	return nil
}
