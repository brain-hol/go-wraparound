package main

import (
	"net/http"

	"github.com/brain-hol/go-httpkit"
	"github.com/brain-hol/go-httpkit/middleware"
)

func main() {
	transport := &httpkit.Transport{}
	transport.Use(middleware.DebugLog)
	client := &http.Client{
		Transport: transport,
	}
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	_, _ = client.Do(req)
}
