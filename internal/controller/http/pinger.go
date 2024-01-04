package httpcontroller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vadimpk/url-pinger/internal/entity"
	"github.com/vadimpk/url-pinger/internal/service"
)

func (c *controller) setupPingerRoutes(rg *gin.RouterGroup) {
	pinger := &pingerController{
		routerContext: routerContext{
			logger:  c.logger.Named("pinger"),
			config:  c.config,
			service: c.service,
		},
	}

	rg.GET("/ping-urls", c.methodWrapper(pinger.pingURLs))
}

type pingerController struct {
	routerContext
}

type pingURLsRequestBody struct {
	URLs        []string `json:"urls" binding:"required"`
	ReturnOnErr bool     `json:"return_on_err"`
	Timeout     int      `json:"timeout"`
}

type pingURLsResponseBody struct {
	Results map[string]entity.URLStatus `json:"results"`
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

	result, err := c.service.PingerService.PingURLs(ctx, service.PingURLOptions{
		URLs:        body.URLs,
		ReturnOnErr: body.ReturnOnErr,
		Timeout:     body.Timeout,
	})
	if err != nil {
		logger.Error("failed to ping urls", "err", err)
		return nil, &httpResponseError{
			Message:    "failed to ping urls",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return pingURLsResponseBody{
		Results: result.Results,
	}, nil
}
