package middleware

import (
	"net/http"

	"github.com/brain-hol/go-httpkit"
)

// Header returns a [RoundTripperWrapper] that adds or overrides a header with
// the specified key and value on each request.
func Header(key string, value string) httpkit.Middleware {
	return httpkit.NewMiddleware(func(req *http.Request, next http.RoundTripper) (*http.Response, error) {
		if key != "" {
			req.Header.Set(key, value)
		}
		return next.RoundTrip(req)
	})
}
