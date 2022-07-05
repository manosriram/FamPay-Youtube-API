package main

import (
	"fmt"
	"log"
	"time"

	"github.com/manosriram/youtubeAPI-fampay/api"
	"github.com/manosriram/youtubeAPI-fampay/db"
	config "github.com/manosriram/youtubeAPI-fampay/pkg/config"
	"github.com/manosriram/youtubeAPI-fampay/pkg/ytservice"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

const PORT = ":5001"

func refreshVideoList(logger *zap.SugaredLogger, config config.Config, collection *mongo.Collection) {
	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				logger.Infow("loading video list")
				ytservice.FetchVideosByQuery(logger, config, "football", collection)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func main() {
	r := gin.New()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("error creating logger", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	config, err := config.Load(".")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(config.YoutubeDeveloperKey)

	collection, err := db.ConnectMongo(sugar)
	if err != nil {
		sugar.Errorw("error connecting to mongo", "error", err)
	} else {
		sugar.Infow("connected to mongo")
	}

	youtubeApi := api.YoutubeAPI{
		Config:          &config,
		MongoCollection: collection,
		Logger:          sugar,
	}

	// go refreshVideoList(sugar, config, collection)

	errorHandler := api.ProvideAPIWrap(sugar)

	r.GET("/", errorHandler(youtubeApi.LoadStoredVideos))

	r.GET("/search", errorHandler(youtubeApi.LoadStoredVideosByQuery))

	r.Run(PORT)
}
