package httpcontroller

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/DataDog/gostackparse"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	requestIDMiddlewareKey = "requestID"
)

func (c *controller) requestIDMiddleware(ctx *gin.Context) {
	ctx.Set(requestIDMiddlewareKey, uuid.NewString())
}

func (c *controller) panicHandlerMiddleware() gin.HandlerFunc {
	logger := c.logger.Named("panicHandler")
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// get stacktrace
				stacktrace, errors := gostackparse.Parse(bytes.NewReader(debug.Stack()))
				if len(errors) > 0 || len(stacktrace) == 0 {
					logger.Error("get stacktrace errors", "stacktraceErrors", errors, "stacktrace", "unknown", "err", err)
				} else {
					logger.Error("unhandled error", "err", err, "stacktrace", stacktrace)
				}
				// return error
				err := c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("%v", err))
				if err != nil {
					logger.Error("failed to abort with error", "err", err)
				}
			}
		}()
		c.Next()
	}
}
