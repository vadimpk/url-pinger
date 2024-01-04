package httpcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vadimpk/url-pinger/config"
	logging "github.com/vadimpk/url-pinger/pkg/logger"
)

type Options struct {
	Logger logging.Logger
	Config *config.Config
}

type routerContext struct {
	logger logging.Logger
	config *config.Config
}

type controller struct {
	routerContext
}

func New(opts Options) http.Handler {
	r := gin.New()

	c := &controller{
		routerContext: routerContext{
			logger: opts.Logger.Named("HTTPController"),
			config: opts.Config,
		},
	}

	rg := r.Group("/api/v1")

	rg.Use(
		c.panicHandlerMiddleware(),
		c.requestIDMiddleware,
	)

	c.setupPingerRoutes(rg)

	return r
}

type httpResponseError struct {
	Message    string `json:"message"`
	StatusCode int
}

func (c *controller) methodWrapper(handler func(ctx *gin.Context) (interface{}, *httpResponseError)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		result, err := handler(ctx)
		if err != nil {
			if err.StatusCode == 0 {
				err.StatusCode = http.StatusInternalServerError
			}

			ctx.AbortWithStatusJSON(err.StatusCode, err)
			return
		}
		ctx.JSON(http.StatusOK, result)
	}
}
