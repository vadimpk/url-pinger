package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlValidator_ValidateURL(t *testing.T) {
	t.Parallel()

	validator := NewURLValidator()

	testCases := []struct {
		name     string
		url      string
		expected bool
	}{
		{
			name:     "valid url",
			url:      "https://google.com",
			expected: true,
		},
		{
			name:     "invalid url",
			url:      "google.com",
			expected: false,
		},
		{
			name:     "invalid url path",
			url:      "https://google.com/!@#$%^&*()",
			expected: false,
		},
		{
			name:     "empty url",
			url:      "",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := validator.ValidateURL(tc.url)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
