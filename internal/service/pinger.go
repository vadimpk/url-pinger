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
	wg      *sync.WaitGroup
	ctx     context.Context
	url     string
	resCh   chan<- result
	timeout time.Duration
}

type result struct {
	url    string
	status entity.URLStatus
	t      time.Duration
}

func (s *pingerService) PingURLs(ctx context.Context, options PingURLOptions) (*PingURLsResponse, error) {
	wg := sync.WaitGroup{}
	wg.Add(len(options.URLs))
	respCh := make(chan result, len(options.URLs))

	var timeout time.Duration
	if options.Timeout <= 0 {
		timeout = s.config.Ping.DefaultPingTimeout
	} else {
		timeout = time.Duration(options.Timeout) * time.Second
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		for _, url := range options.URLs {
			s.jobs <- job{
				wg:      &wg,
				ctx:     ctx,
				url:     url,
				resCh:   respCh,
				timeout: timeout,
			}
		}
	}()

	go func() {
		wg.Wait()
		close(respCh)
	}()

	results := make(map[string]entity.URLStatus)
	totalTime := time.Duration(0)
	for res := range respCh {
		results[res.url] = res.status
		totalTime += res.t

		if options.ReturnOnErr && res.status != entity.URLStatusOK {
			cancel()
			break
		}
	}

	return &PingURLsResponse{
		Results:           results,
		AverageTimePerURL: totalTime / time.Duration(len(options.URLs)),
	}, nil
}

func (s *pingerService) worker() {
	for j := range s.jobs {
		select {
		case <-j.ctx.Done():
			j.wg.Done()
			continue
		default:
			j.resCh <- s.do(j)
			j.wg.Done()
			continue
		}
	}
}

func (s *pingerService) do(j job) result {
	start := time.Now()
	valid := s.urlValidator.ValidateURL(j.url)
	if !valid {
		s.logger.Info("invalid url", "url", j.url)
		return result{
			url:    j.url,
			status: entity.URLStatusInvalid,
			t:      time.Since(start),
		}
	}

	status, err := s.basePinger.Ping(j.url, PingTimeoutOption{Timeout: j.timeout})
	if err != nil {
		s.logger.Error("failed to ping url", "url", j.url, "err", err)

		return result{
			url:    j.url,
			status: entity.URLStatusError,
			t:      time.Since(start),
		}
	}

	return result{
		url:    j.url,
		status: status,
		t:      time.Since(start),
	}
}
