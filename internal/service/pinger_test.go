package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/vadimpk/url-pinger/config"
	"github.com/vadimpk/url-pinger/internal/entity"
	logging "github.com/vadimpk/url-pinger/pkg/logger"
)

func TestPingerService_PingURLs(t *testing.T) {
	t.Parallel()

	baseOptions := BaseOptions{
		Logger: logging.New("debug"),
		Config: config.Get(),
	}

	testCases := []struct {
		name         string
		opts         PingURLOptions
		mockedPinger URLPinger
		baseOptions  BaseOptions
		expected     *PingURLsResponse
	}{
		{
			name: "ping urls: all ok",
			opts: PingURLOptions{
				URLs: []string{
					"https://google.com",
					"http://www.wikipedia.org",
					"http://example.com",
				},
				ReturnOnErr: false,
			},
			mockedPinger: NewPingerMock(nil, entity.URLStatusOK),
			baseOptions:  baseOptions,
			expected: &PingURLsResponse{
				Results: map[string]entity.URLStatus{
					"https://google.com":       entity.URLStatusOK,
					"http://www.wikipedia.org": entity.URLStatusOK,
					"http://example.com":       entity.URLStatusOK,
				},
			},
		},
		{
			name: "ping urls: one failed",
			opts: PingURLOptions{
				URLs: []string{
					"https://google.com",
					"http://www.wikipedia.org",
					"http://example.com",
				},
				ReturnOnErr: false,
			},
			mockedPinger: NewPingerMock(nil, entity.URLStatusOK).On("http://example.com", nil, entity.URLStatusFailed),
			baseOptions:  baseOptions,
			expected: &PingURLsResponse{
				Results: map[string]entity.URLStatus{
					"https://google.com":       entity.URLStatusOK,
					"http://www.wikipedia.org": entity.URLStatusOK,
					"http://example.com":       entity.URLStatusFailed,
				},
			},
		},
		{
			name: "ping urls: first failed, return on err",
			opts: PingURLOptions{
				URLs: []string{
					"https://google.com",
					"http://www.wikipedia.org",
					"http://example.com",
				},
				ReturnOnErr: true,
			},
			mockedPinger: NewPingerMock(nil, entity.URLStatusOK).On("https://google.com", nil, entity.URLStatusFailed),
			baseOptions: BaseOptions{
				Logger: baseOptions.Logger,
				Config: &config.Config{
					Ping: config.PingerService{
						WorkerPoolSize: 1,
					},
				},
			},
			expected: &PingURLsResponse{
				Results: map[string]entity.URLStatus{
					"https://google.com": entity.URLStatusFailed,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			service := NewPingerService(tc.baseOptions, tc.mockedPinger, NewURLValidator())

			actual, err := service.PingURLs(context.Background(), tc.opts)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if !reflect.DeepEqual(tc.expected.Results, actual.Results) {
				t.Errorf("expected results: %v, actual: %v", tc.expected.Results, actual.Results)
			}
		})
	}
}

type pingerMock struct {
	defaultResponse pingerMockResponse
	responses       map[string]pingerMockResponse
}

type pingerMockResponse struct {
	err    error
	status entity.URLStatus
}

func NewPingerMock(err error, status entity.URLStatus) *pingerMock {
	return &pingerMock{
		defaultResponse: pingerMockResponse{
			err:    err,
			status: status,
		},
		responses: make(map[string]pingerMockResponse),
	}
}

func (p *pingerMock) On(url string, err error, status entity.URLStatus) *pingerMock {
	p.responses[url] = pingerMockResponse{
		err:    err,
		status: status,
	}

	return p
}

func (p *pingerMock) Ping(url string, opts ...interface{}) (entity.URLStatus, error) {
	if resp, ok := p.responses[url]; ok {
		return resp.status, resp.err
	}

	return p.defaultResponse.status, p.defaultResponse.err
}

func (p *pingerMock) SetNext(next URLPinger) {
}
