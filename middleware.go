package httpkit

import "net/http"

// Transport -------------------------------------------------------------------

// Transport is a middleware-enabled wrapper around an [http.RoundTripper] that
// executes a chain of [Middleware] added via the Use method.
type Transport struct {
	// Transport is the underlying [http.RoundTripper]. If nil,
	// [http.DefaultTransport] is used.
	Transport http.RoundTripper
}

// RoundTrip executes the request using the [Transport]'s middleware chain and
// the underlying transport. It implements the [http.RoundTripper] interface.
func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.Transport == nil {
		return http.DefaultTransport.RoundTrip(req)
	}
	return t.Transport.RoundTrip(req)
}

// Use appends one or more [Middleware] to the [Transport].
func (t *Transport) Use(middlewares ...Middleware) {
	base := t.Transport
	if base == nil {
		base = http.DefaultTransport
	}

	for i := len(middlewares) - 1; i >= 0; i-- {
		base = middlewares[i](base)
	}

	t.Transport = base
}

// Middleware ------------------------------------------------------------------

// Middleware defines a function that wraps an [http.RoundTripper], allowing
// custom middleware behavior to be injected into the request lifecycle.
type Middleware func(next http.RoundTripper) http.RoundTripper

// NewMiddleware converts a HandlerFunc into a Middleware.
// This helper lets users write middleware with the 'next' parameter clearly named.
func NewMiddleware(h func(req *http.Request, next http.RoundTripper) (*http.Response, error)) Middleware {
	return func(next http.RoundTripper) http.RoundTripper {
		return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return h(req, next)
		})
	}
}

// roundTripperFunc is an unexported adapter to allow ordinary functions as http.RoundTripper.
type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
