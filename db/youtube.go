package db

import (
	"context"
	"fmt"
	"time"

	"github.com/manosriram/youtubeAPI-fampay/data"
	"github.com/manosriram/youtubeAPI-fampay/pkg/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

const (
	VideosPerPage  = 5
	DatabaseName   = "youtubeapi"
	CollectionName = "youtubeapi"
)

/*
   connects to mongodb cloud and returns the collection
*/
func ConnectMongo(logger *zap.SugaredLogger, config config.Config) (*mongo.Collection, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)

	mongoURI := fmt.Sprintf("mongodb+srv://%s:%s@youtube-api-fampay.g7ssryz.mongodb.net/?retryWrites=true&w=majority", config.MongoUsername, config.MongoPassword)
	fmt.Println(mongoURI)

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)

	collection := client.Database(DatabaseName).Collection(CollectionName)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

/*
   bulk upserts list of videos, inserts if video doesn't exist; updates if existing video has any changes in it's fields
*/
func BulkUpsertVideos(ctx context.Context, logger *zap.SugaredLogger, videos []data.Video, mongoCollection *mongo.Collection) error {
	var operations []mongo.WriteModel

	for _, video := range videos {
		operation := mongo.NewUpdateOneModel()
		operation.SetUpsert(true)
		operation.SetFilter(bson.M{"video_etag": video.VideoETag})
		operation.SetUpdate(bson.M{"$set": bson.M{"title": video.Title, "description": video.Description, "thumbnail_url": video.ThumbnailUrl, "published_at": video.PublishedAt}})

		operations = append(operations, operation)
	}

	bulkWrite := options.BulkWriteOptions{}
	bulkWrite.SetOrdered(true)

	_, err := mongoCollection.BulkWrite(ctx, operations, &bulkWrite)
	if err != nil {
		logger.Errorw("error calling BulkWrite", "error", err)
		return err
	}
	return nil
}

/*
   fetches video list from mongodb cloud.
   supports pagination and search query.
*/
func GetVideosList(ctx context.Context, showVideoRequest data.ShowVideoRequest, logger *zap.SugaredLogger, mongoCollection *mongo.Collection) ([]*data.Video, error) {
	videoList := make([]*data.Video, 0)

	var (
		page        = showVideoRequest.Page
		searchQuery = showVideoRequest.SearchQuery
	)

	filter := bson.M{}

	// if searchQuery isn't empty, do partial search on title and description
	if searchQuery != "" {
		filter = bson.M{"$and": []bson.M{
			{"$text": bson.M{"$search": searchQuery}},
		},
		}
	}

	// sort in descending order of published_at field
	findOptions := &options.FindOptions{}
	findOptions.SetSort(bson.M{"published_at": -1})

	// skip documents to reach given page, considering 5 videos per page
	findOptions.SetSkip(int64((page - 1) * VideosPerPage))
	findOptions.SetLimit(4)

	cursor, err := mongoCollection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		video := &data.Video{}
		err := cursor.Decode(video)
		if err != nil {
			return nil, err
		} else {
			videoList = append(videoList, video)
		}
	}

	return videoList, nil
}
