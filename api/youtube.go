package api

import (
	"net/http"
	"strconv"

	"github.com/manosriram/youtubeAPI-fampay/data"
	config "github.com/manosriram/youtubeAPI-fampay/pkg/config"
	ytservice "github.com/manosriram/youtubeAPI-fampay/pkg/ytservice"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type YoutubeAPI struct {
	Logger          *zap.SugaredLogger
	Config          *config.Config
	YoutubeService  ytservice.Service
	MongoCollection *mongo.Collection
}

func NewYouTubeAPI(logger *zap.SugaredLogger, config *config.Config, youtubeService ytservice.Service, mongoCollection *mongo.Collection) *YoutubeAPI {
	return &YoutubeAPI{
		Logger:          logger,
		Config:          config,
		YoutubeService:  youtubeService,
		MongoCollection: mongoCollection,
	}
}

/*
   method: GET
   params: page
   description: loads video(s) metadata from db and returns an json array

*/
func (yt YoutubeAPI) LoadStoredVideos(ctx *gin.Context) (int, interface{}, error) {
	logger := yt.Logger

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		logger.Errorw("error getting page query", "error", err)
		return http.StatusInternalServerError, gin.H{"success": false, "error": "error getting page"}, err
	}

	showVideoRequest := data.ShowVideoRequest{
		Page:        int64(page),
		SearchQuery: "",
	}

	videos, err := yt.YoutubeService.LoadStoredVideos(ctx, showVideoRequest, yt.Logger, yt.MongoCollection)
	if err != nil {
		logger.Errorw("error getting videos list", "error", err)
		return http.StatusInternalServerError, gin.H{"success": false, "error": "error loading videos"}, err
	}

	return http.StatusOK, gin.H{"success": true, "videos": videos, "error": ""}, nil
}

/*
    method: GET
    params: page, search
    description: loads video(s) matching the query metadata from db and returns an json array.
		 supports partial matching on fields - title and description

*/
func (yt YoutubeAPI) LoadStoredVideosByQuery(ctx *gin.Context) (int, interface{}, error) {
	logger := yt.Logger

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		logger.Errorw("error getting page query", "error", err)
		return http.StatusInternalServerError, gin.H{"success": false, "error": "error getting page"}, err
	}

	searchQuery := ctx.Query("search")

	showVideoRequest := data.ShowVideoRequest{
		Page:        int64(page),
		SearchQuery: searchQuery,
	}

	videos, err := yt.YoutubeService.LoadStoredVideos(ctx, showVideoRequest, yt.Logger, yt.MongoCollection)
	if err != nil {
		logger.Errorw("error getting videos list", "error", err)
		return http.StatusInternalServerError, gin.H{"success": false, "error": "error loading videos"}, err
	}

	return http.StatusOK, gin.H{"success": true, "videos": videos, "error": ""}, nil
}
