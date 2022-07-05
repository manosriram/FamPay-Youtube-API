package main

import (
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

// function which calls FetchVideosByQuery(...) every 10seconds
func refreshVideoList(logger *zap.SugaredLogger, config *config.Config, collection *mongo.Collection) {
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

	// initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("error creating logger", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	// initialize config, get variables from config.yaml
	config, err := config.Load(".")
	if err != nil {
		log.Fatal(err)
	}

	// connect to mongodb service (cloud)
	collection, err := db.ConnectMongo(sugar, config)
	if err != nil {
		sugar.Errorw("error connecting to mongo", "error", err)
	} else {
		sugar.Infow("connected to mongo")
	}

	// create youtube service
	youtubeService := ytservice.Service{
		Config: &config,
	}

	// create youtube-api to access APIs
	youtubeApi := api.YoutubeAPI{
		Config:          &config,
		MongoCollection: collection,
		YoutubeService:  youtubeService,
		Logger:          sugar,
	}

	// go refreshVideoList(sugar, &config, collection)

	// error handler which intercepts response
	errorHandler := api.ProvideAPIWrap(sugar)

	// handler functions
	r.GET("/", errorHandler(youtubeApi.LoadStoredVideos))

	r.GET("/search", errorHandler(youtubeApi.LoadStoredVideosByQuery))

	// start the server
	r.Run(PORT)
}
