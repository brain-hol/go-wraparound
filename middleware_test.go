package httpkit

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

func TestTransport_CorrectOrder(t *testing.T) {
	assert := assert.New(t)

	callSequence := []string{}

	middlewareA := NewMiddleware(func(req *http.Request, next http.RoundTripper) (*http.Response, error) {
		callSequence = append(callSequence, "A")
		resp, err := next.RoundTrip(req)
		callSequence = append(callSequence, "D")
		return resp, err
	})

	middlewareB := NewMiddleware(func(req *http.Request, next http.RoundTripper) (*http.Response, error) {
		callSequence = append(callSequence, "B")
		resp, err := next.RoundTrip(req)
		callSequence = append(callSequence, "C")
		return resp, err
	})

	w := &Transport{
		Transport: &mockRoundTripper{
			Response: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("wrapped")),
			},
		},
	}

	w.Use(middlewareA, middlewareB)

	req, err := http.NewRequest("GET", "http://example.com", nil)
	assert.NoError(err)

	_, err = w.RoundTrip(req)
	assert.NoError(err)
	assert.Equal([]string{"A", "B", "C", "D"}, callSequence)
}
