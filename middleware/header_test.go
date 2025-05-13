package middleware

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockRoundTripper struct {
	Response *http.Response
	Err      error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

func TestHeaderMiddleware(t *testing.T) {
	tests := []struct {
		name        string
		headerKey   string
		headerValue string
	}{
		{"Well-known header works", "Authorization", "Bearer example-token"},
		{"Unknown header works", "Random", "random value"},
		{"No value set if key is empty", "", "random value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			mockResp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("OK")),
			}
			mockRT := &mockRoundTripper{
				Response: mockResp,
			}

			wrapper := Header(tt.headerKey, tt.headerValue)(mockRT)

			req, err := http.NewRequest("GET", "http://example.com", nil)
			assert.NoError(err)

			resp, err := wrapper.RoundTrip(req)
			assert.NoError(err)
			assert.Equal(http.StatusOK, resp.StatusCode)

			got := req.Header.Get(tt.headerKey)
			if tt.headerKey != "" {
				assert.Equal(tt.headerValue, got)
			} else {
				assert.Equal("", got)
			}
		})
	}
}
