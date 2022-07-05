package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller func(*gin.Context) (int, interface{}, error)
type Wrapper func(Controller) gin.HandlerFunc

func ProvideAPIWrap(logger *zap.SugaredLogger) Wrapper {
	return func(controller Controller) gin.HandlerFunc {
		return func(ctx *gin.Context) {
			status, body, err := controller(ctx)
			if err != nil {

			}

			ctx.JSON(status, body)
		}
	}
}
