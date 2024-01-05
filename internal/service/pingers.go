package service

import (
	"context"
	"errors"
	"net/http"

	"github.com/vadimpk/url-pinger/internal/entity"
)

type httpPinger struct {
	serviceContext
	client *http.Client
	Next   URLPinger
}

func NewHTTPPinger(options BaseOptions) URLPinger {
	return &httpPinger{
		serviceContext: serviceContext{
			logger: options.Logger.Named("HTTPPinger"),
			config: options.Config,
		},
		client: &http.Client{},
	}
}

func (p *httpPinger) Ping(url string, opts ...interface{}) (entity.URLStatus, error) {
	ctx := context.Background()
	for _, opt := range opts {
		switch opt.(type) {
		case PingTimeoutOption:
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, opt.(PingTimeoutOption).Timeout)
			defer cancel()
		}
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return entity.URLStatusFailed, err
	}

	resp, err := p.client.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return entity.URLStatusTimeout, nil
		}
		return entity.URLStatusError, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return entity.URLStatusNotFound, nil
	}

	if resp.StatusCode != http.StatusOK {
		return entity.URLStatusError, nil
	}

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

func (p *cachePinger) Ping(url string, opts ...interface{}) (entity.URLStatus, error) {
	// TODO: add cache logic (for now just call the next pinger)

	if p.Next != nil {
		return p.Next.Ping(url, opts...)
	}

	return entity.URLStatusUnknown, nil
}

func (p *cachePinger) SetNext(next URLPinger) {
	p.Next = next
}
