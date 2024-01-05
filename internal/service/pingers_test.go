package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vadimpk/url-pinger/config"
	"github.com/vadimpk/url-pinger/internal/entity"
	logging "github.com/vadimpk/url-pinger/pkg/logger"
)

func TestHttpPinger_Ping(t *testing.T) {
	t.Parallel()

	pinger := NewHTTPPinger(BaseOptions{
		Logger: logging.New("debug"),
		Config: config.Get(),
	})

	testCases := []struct {
		name     string
		url      string
		expected entity.URLStatus
	}{
		{
			name:     "ok url",
			url:      "https://google.com",
			expected: entity.URLStatusOK,
		},
		{
			name:     "not found url",
			url:      "https://example.com/404",
			expected: entity.URLStatusNotFound,
		},
		{
			name:     "error url",
			url:      "localhost:1234",
			expected: entity.URLStatusError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, _ := pinger.Ping(tc.url)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
