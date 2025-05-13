package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/brain-hol/go-httpkit"
)

var DebugLog = httpkit.NewMiddleware(func(req *http.Request, next http.RoundTripper) (*http.Response, error) {
	fmt.Println(gray + "╭───────────────────────────────────────────────────────────────────────────────" + reset)

	// Log Request
	fmt.Printf(gray+"│ ▶ REQUEST: "+methodColor(req.Method)+"%s"+reset+" %s\n", req.Method, req.URL)
	for k, v := range req.Header {
		fmt.Printf(gray+"│   "+cyan+"%s"+reset+": %s\n", k, strings.Join(v, ", "))
	}

	var reqBody []byte
	if req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewBuffer(reqBody)) // Restore body
		logBody(reqBody, req.Header.Get("Content-Type"))
	}
	fmt.Println(gray + "│")

	// Forward request
	resp, err := next.RoundTrip(req)
	if err != nil {
		fmt.Println(gray + "│   ERROR: " + red + err.Error() + reset)
		fmt.Println(gray + "╰───────────────────────────────────────────────────────────────────────────────" + reset)
		return nil, err
	}

	// Log Response
	fmt.Printf(gray+"│ ◀ RESPONSE: "+statusColor(resp.StatusCode)+"%s\n", resp.Status)
	for k, v := range resp.Header {
		fmt.Printf(gray+"│   "+cyan+"%s"+reset+": %s\n", k, strings.Join(v, ", "))
	}

	var respBody []byte
	if resp.Body != nil {
		respBody, _ = io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody)) // Restore body
		logBody(respBody, resp.Header.Get("Content-Type"))
	}

	fmt.Println(gray + "╰───────────────────────────────────────────────────────────────────────────────" + reset)
	return resp, nil
})

const (
	gray   = "\033[90m"
	cyan   = "\033[36m"
	green  = "\033[32m"
	blue   = "\033[34m"
	yellow = "\033[33m"
	red    = "\033[31m"
	reset  = "\033[0m"
)

func methodColor(method string) string {
	switch strings.ToUpper(method) {
	case "GET":
		return green
	case "POST":
		return blue
	case "PUT", "PATCH":
		return yellow
	case "DELETE":
		return red
	default:
		return reset
	}
}

func statusColor(status int) string {
	switch {
	case status >= 200 && status < 300:
		return green
	case status >= 300 && status < 400:
		return yellow
	case status >= 400:
		return red
	default:
		return reset
	}
}

func getBodyLines(body []byte, contentType string) []string {
	if strings.Contains(contentType, "application/json") {
		var out bytes.Buffer
		if err := json.Indent(&out, body, "", "  "); err == nil {
			return strings.Split(out.String(), "\n")
		}
	}
	return strings.Split(string(body), "\n")
}

func logBody(body []byte, contentType string) {
	if len(body) == 0 || (!strings.Contains(contentType, "application/json") && !strings.HasPrefix(contentType, "text/")) {
		return
	}
	fmt.Println(gray + "│" + reset)
	for _, line := range getBodyLines(body, contentType) {
		fmt.Printf(gray+"│   "+reset+"%s\n", line)
	}
}
