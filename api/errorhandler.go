package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller func(*gin.Context) (int, interface{}, error)
type Wrapper func(Controller) gin.HandlerFunc

// error handler
func ProvideAPIWrap(logger *zap.SugaredLogger) Wrapper {
	return func(controller Controller) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			status, body, err := controller(ctx)
			if err != nil {
				logger.Errorw("error in handler", "error", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": err.Error()})
				return
			}

			ctx.JSON(status, body)
		}
	}
}
