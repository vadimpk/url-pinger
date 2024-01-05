package service

import (
	"context"
	"time"

	"github.com/vadimpk/url-pinger/config"
	"github.com/vadimpk/url-pinger/internal/entity"
	logging "github.com/vadimpk/url-pinger/pkg/logger"
)

type Services struct {
	PingerService
}

type BaseOptions struct {
	Logger logging.Logger
	Config *config.Config
}

type serviceContext struct {
	logger logging.Logger
	config *config.Config
}

type PingerService interface {
	PingURLs(ctx context.Context, options PingURLOptions) (*PingURLsResponse, error)
}

type PingURLOptions struct {
	URLs        []string
	ReturnOnErr bool
	Timeout     int
}

type PingURLsResponse struct {
	Results map[string]entity.URLStatus
}

type URLValidator interface {
	ValidateURL(url string) (bool, error)
}

type URLPinger interface {
	Ping(ctx context.Context, url string, opts ...interface{}) (entity.URLStatus, error)
	SetNext(next URLPinger)
}

type PingTimeoutOption struct {
	Timeout time.Duration
}

type CacheStorage interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}
