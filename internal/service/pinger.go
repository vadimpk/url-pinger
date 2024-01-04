package service

import (
	"context"
)

type pingerService struct {
	serviceContext
	basePinger   URLPinger
	urlValidator URLValidator
}

func NewPingerService(options BaseOptions, basePinger URLPinger, validator URLValidator) PingerService {
	return &pingerService{
		serviceContext: serviceContext{
			logger: options.Logger.Named("PingerService"),
			config: options.Config,
		},
		basePinger:   basePinger,
		urlValidator: validator,
	}
}

func (s *pingerService) PingURLs(ctx context.Context, options PingURLOptions) (*PingURLsResponse, error) {
	return nil, nil
}
