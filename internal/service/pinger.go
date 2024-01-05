package service

import (
	"context"
	"sync"
	"time"

	"github.com/vadimpk/url-pinger/internal/entity"
)

type pingerService struct {
	serviceContext
	basePinger   URLPinger
	urlValidator URLValidator

	jobs chan job
}

func NewPingerService(options BaseOptions, basePinger URLPinger, validator URLValidator) PingerService {
	s := &pingerService{
		serviceContext: serviceContext{
			logger: options.Logger.Named("PingerService"),
			config: options.Config,
		},
		basePinger:   basePinger,
		urlValidator: validator,
		jobs:         make(chan job, options.Config.Ping.WorkerPoolSize),
	}

	for i := 0; i < s.config.Ping.WorkerPoolSize; i++ {
		go s.worker()
	}

	return s
}

type job struct {
	ctx     context.Context
	url     string
	resCh   chan<- result
	stop    <-chan struct{}
	timeout time.Duration
}

type result struct {
	url    string
	status entity.URLStatus
}

func (s *pingerService) PingURLs(ctx context.Context, options PingURLOptions) (*PingURLsResponse, error) {
	wg := sync.WaitGroup{}
	wg.Add(len(options.URLs))
	respCh := make(chan result)
	stop := make(chan struct{})

	var timeout time.Duration
	if options.Timeout <= 0 {
		timeout = s.config.Ping.DefaultPingTimeout
	} else {
		timeout = time.Duration(options.Timeout) * time.Second
	}

	go func() {
		for _, url := range options.URLs {
			s.jobs <- job{
				ctx:     ctx,
				url:     url,
				resCh:   respCh,
				stop:    stop,
				timeout: timeout,
			}
		}
	}()

	go func() {
		wg.Wait()
		close(respCh)
	}()

	results := make(map[string]entity.URLStatus)
	for res := range respCh {
		results[res.url] = res.status

		if options.ReturnOnErr && res.status != entity.URLStatusOK {
			close(stop)
		}

		wg.Done()
	}

	return &PingURLsResponse{
		Results: results,
	}, nil
}

func (s *pingerService) worker() {
	for j := range s.jobs {
		select {
		case <-j.stop:
			return
		case <-j.ctx.Done():
			return
		default:
			valid, err := s.urlValidator.ValidateURL(j.url)
			if err != nil {
				s.logger.Error("failed to validate url", "url", j.url, "err", err)

				j.resCh <- result{
					url:    j.url,
					status: entity.URLStatusFailed,
				}

				continue
			}

			if !valid {
				j.resCh <- result{
					url:    j.url,
					status: entity.URLStatusInvalid,
				}

				continue
			}

			status, err := s.basePinger.Ping(j.ctx, j.url, PingTimeoutOption{Timeout: j.timeout})
			if err != nil {
				s.logger.Error("failed to ping url", "url", j.url, "err", err)

				j.resCh <- result{
					url:    j.url,
					status: entity.URLStatusFailed,
				}

				continue
			}

			j.resCh <- result{
				url:    j.url,
				status: status,
			}

			continue
		}
	}
}
