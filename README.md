# httpkit

Middleware-enabled wrapper around `http.RoundTripper`.

`httpkit` is a Go package that provides a flexible way to chain and apply custom middleware around an [http.RoundTripper](https://pkg.go.dev/net/http#RoundTripper) (a.k.a. transport) to modify outbound HTTP requests and their respective responses. Think of it as middleware for the `http.Client` instead of the `http.Server`.

The package wraps the standard library's `http.RoundTripper`, letting you easily create and use custom middleware to inject behavior into the HTTP client request lifecycle — all without external dependencies.

This package is inspired by HTTP server middleware patterns found in the Go standard library and frameworks like [chi](https://github.com/go-chi/chi).

## Installation

Install the package using:

```shell
go get github.com/brain-hol/go-httpkit
```

Then import it in your Go code:

```go
import "github.com/brain-hol/go-httpkit"
```

## Usage

### Basic Example

```go
package main

import (
	"fmt"
	"net/http"
	"github.com/brain-hol/go-httpkit"
)

func main() {
	// Create a new httpkit.Transport
	transport := &httpkit.Transport{}

	// Add a middleware that adds a custom header
	transport.Use(httpkit.NewMiddleware(func(req *http.Request, next http.RoundTripper) (*http.Response, error) {
		req.Header.Set("Authorization", "Bearer my-token")
		return next.RoundTrip(req)
	}))

	// Create an HTTP client that uses the httpkit.Transport
	client := &http.Client{
		Transport: transport,
	}

	// Send a request through the middleware chain
	req, _ := http.NewRequest("GET", "https://example.com/resource", nil)
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Response Status:", resp.Status)
	}
}
```

In this example, the middleware automatically adds an authorization header to every request.

### Adding Custom Middleware

You can define reusable middleware variables or functions. For example:

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/brain-hol/go-httpkit"
)

// Logging middleware as a variable
var Logging = httpkit.NewMiddleware(func(req *http.Request, next http.RoundTripper) (*http.Response, error) {
	log.Printf("Sending request %s %s", req.Method, req.URL)
	resp, err := next.RoundTrip(req)
	if err != nil {
		log.Printf("Request error: %v", err)
	} else {
		log.Printf("Received response status: %s", resp.Status)
	}
	return resp, err
})

func main() {
	transport := &httpkit.Transport{}
	transport.Use(Logging)

	client := &http.Client{Transport: transport}

	req, _ := http.NewRequest("GET", "https://example.com", nil)
	client.Do(req)
}
```

This shows how to declare middleware as a variable using `NewMiddleware` for clean and simple composition.

If you want middleware that needs configuration, you can write middleware **constructor functions** that return `httpkit.Middleware` with parameters.
