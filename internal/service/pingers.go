package service

import (
	"context"

	"github.com/vadimpk/url-pinger/internal/entity"
)

type httpPinger struct {
	serviceContext
	Next URLPinger
}

func NewHTTPPinger(options BaseOptions) URLPinger {
	return &httpPinger{
		serviceContext: serviceContext{
			logger: options.Logger.Named("HTTPPinger"),
			config: options.Config,
		},
	}
}

func (p *httpPinger) Ping(ctx context.Context, url string) (entity.URLStatus, error) {
	return entity.URLStatusOK, nil
}

func (p *httpPinger) SetNext(next URLPinger) {
	p.Next = next
}

type cachePinger struct {
	serviceContext
	storage CacheStorage
	Next    URLPinger
}

func NewCachePinger(options BaseOptions, storage CacheStorage) URLPinger {
	return &cachePinger{
		serviceContext: serviceContext{
			logger: options.Logger.Named("CachePinger"),
			config: options.Config,
		},
		storage: storage,
	}
}

func (p *cachePinger) Ping(ctx context.Context, url string) (entity.URLStatus, error) {
	// TODO: add cache logic (for now just call the next pinger)

	if p.Next != nil {
		return p.Next.Ping(ctx, url)
	}

	return entity.URLStatusUnknown, nil
}

func (p *cachePinger) SetNext(next URLPinger) {
	p.Next = next
}
