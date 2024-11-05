package polycode

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"mime"
	"net/http"
	"strings"
)

// ResponseWriter implements the http.ResponseWriter interface
// in order to support the API Gateway Lambda HTTP "protocol".
type ResponseWriter struct {
	out           ApiResponse
	buf           bytes.Buffer
	header        http.Header
	wroteHeader   bool
	closeNotifyCh chan bool
}

// Header implementation.
func (w *ResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}

	return w.header
}

// Write implementation.
func (w *ResponseWriter) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}

	return w.buf.Write(b)
}

// WriteHeader implementation.
func (w *ResponseWriter) WriteHeader(status int) {
	if w.wroteHeader {
		return
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf8")
	}

	w.out.StatusCode = status

	h := make(map[string]string)
	for k, v := range w.Header() {
		if len(v) == 1 {
			h[k] = v[0]
		} else if len(v) > 1 {
			h[k] = strings.Join(v, ",")
		}
	}

	w.out.Header = h
	w.wroteHeader = true
}

// CloseNotify notify when the response is closed
func (w *ResponseWriter) CloseNotify() <-chan bool {
	return w.closeNotifyCh
}

// End the request.
func (w *ResponseWriter) End() ApiResponse {
	w.out.IsBase64Encoded = isBinary(w.header)

	if w.out.IsBase64Encoded {
		w.out.Body = base64.StdEncoding.EncodeToString(w.buf.Bytes())
	} else {
		w.out.Body = w.buf.String()
	}

	// notify end
	//w.closeNotifyCh <- true

	return w.out
}

// isBinary returns true if the response reprensents binary.
func isBinary(h http.Header) bool {
	switch {
	case !isTextMime(h.Get("Content-Type")):
		return true
	case h.Get("Content-Encoding") == "gzip":
		return true
	default:
		return false
	}
}

// isTextMime returns true if the content type represents textual data.
func isTextMime(kind string) bool {
	mt, _, err := mime.ParseMediaType(kind)
	if err != nil {
		return false
	}

	if strings.HasPrefix(mt, "text/") {
		return true
	}

	switch mt {
	case "image/svg+xml", "application/json", "application/xml", "application/javascript":
		return true
	default:
		return false
	}
}

type SerializableRequest struct {
	Method string
	URL    string
	Header http.Header
	Body   []byte
}

func convertToHttpRequest(ctx context.Context, apiReq ApiRequest) (*http.Request, error) {
	// Build the URL
	url := apiReq.Path
	if len(apiReq.Query) > 0 {
		queryParams := "?"
		for key, value := range apiReq.Query {
			queryParams += key + "=" + value + "&"
		}
		queryParams = strings.TrimSuffix(queryParams, "&")
		url += queryParams
	}

	// Create a new HTTP request
	var body io.Reader
	if apiReq.Body != "" {
		body = bytes.NewReader([]byte(apiReq.Body))
	} else {
		body = nil
	}

	println("client: create http request with workflow context")
	req, err := http.NewRequestWithContext(ctx, apiReq.Method, url, body)
	if err != nil {
		return nil, err
	}

	// Add headers
	for key, value := range apiReq.Header {
		req.Header.Set(key, value)
	}

	return req, nil
}

// Manually invoke an HTTP handler
func invokeHandler(handler http.Handler, httpReq *http.Request) ApiResponse {
	customWriter := &ResponseWriter{}
	handler.ServeHTTP(customWriter, httpReq)
	return customWriter.End()
}
