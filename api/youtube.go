package api

import (
	"net/http"

	config "github.com/manosriram/youtubeAPI-fampay/pkg/config"
	ytservice "github.com/manosriram/youtubeAPI-fampay/pkg/ytservice"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type YoutubeAPI struct {
	Logger          *zap.SugaredLogger
	Config          *config.Config
	youtubeService  ytservice.Service
	MongoCollection *mongo.Collection
}

func (yt YoutubeAPI) LoadStoredVideos(ctx *gin.Context) (int, interface{}, error) {
	logger := yt.Logger

	videos, err := yt.youtubeService.LoadStoredVideos(ctx, yt.Logger, yt.MongoCollection)
	if err != nil {
		logger.Errorw("error getting videos list", "error", err)
		return http.StatusInternalServerError, gin.H{"success": false}, err
	}

	return http.StatusOK, videos, nil
}

func (yt YoutubeAPI) LoadStoredVideosByQuery(ctx *gin.Context) (int, interface{}, error) {
	return 0, nil, nil
}
