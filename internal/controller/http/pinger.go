package httpcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (c *controller) setupPingerRoutes(rg *gin.RouterGroup) {
	pinger := &pingerController{
		routerContext: routerContext{
			logger: c.logger.Named("pinger"),
			config: c.config,
		},
	}

	rg.GET("/ping-urls", c.methodWrapper(pinger.pingURLs))
}

type pingerController struct {
	routerContext
}

type pingURLsRequestBody struct {
	URLs []string `json:"urls" binding:"required"`
}

func (c *pingerController) pingURLs(ctx *gin.Context) (interface{}, *httpResponseError) {
	logger := c.logger.Named("ping")

	var body pingURLsRequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.Error("failed to bind json", "err", err)
		return nil, &httpResponseError{
			Message:    "invalid request body",
			StatusCode: http.StatusBadRequest,
		}
	}

	return nil, nil
}
