package db

import (
	"context"
	"time"

	"github.com/manosriram/youtubeAPI-fampay/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func ConnectMongo(logger *zap.SugaredLogger) (*mongo.Collection, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI("mongodb+srv://manosriram:Mano1234$@youtube-api-fampay.g7ssryz.mongodb.net/?retryWrites=true&w=majority").
		SetServerAPIOptions(serverAPIOptions)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	collection := client.Database("youtubeapi").Collection("youtubeapi")
	if err != nil {
		return nil, err
	}
	return collection, nil
}

func BulkUpsertVideos(ctx context.Context, logger *zap.SugaredLogger, videos []data.Video, mongoCollection *mongo.Collection) error {
	var operations []mongo.WriteModel

	for _, video := range videos {
		operation := mongo.NewUpdateOneModel()
		operation.SetUpsert(true)
		operation.SetFilter(bson.M{"video_etag": video.VideoETag})
		operation.SetUpdate(bson.M{"$set": bson.M{"title": video.Title, "description": video.Description, "thumbnail_url": video.ThumbnailUrl}})

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

func GetVideosList(ctx context.Context, logger *zap.SugaredLogger, mongoCollection *mongo.Collection) ([]*data.Video, error) {
	videoList := make([]*data.Video, 0)

	filter := bson.M{}

	cursor, err := mongoCollection.Find(ctx, filter)
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
