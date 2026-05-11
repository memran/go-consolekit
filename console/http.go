package console

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HttpClient struct {
	baseURL     string
	headers     map[string]string
	query       map[string]string
	timeout     time.Duration
	retryCount  int
	retryDelay  time.Duration
	contentType string
	token       string
	username    string
	password    string
}

type Response struct {
	statusCode int
	body       string
	headers    map[string]string
}

func Http() *HttpClient {
	return &HttpClient{
		headers: make(map[string]string),
		query:   make(map[string]string),
	}
}

func (h *HttpClient) BaseURL(url string) *HttpClient {
	h.baseURL = strings.TrimRight(url, "/")
	return h
}

func (h *HttpClient) WithHeader(key, value string) *HttpClient {
	h.headers[key] = value
	return h
}

func (h *HttpClient) WithHeaders(headers map[string]string) *HttpClient {
	for k, v := range headers {
		h.headers[k] = v
	}
	return h
}

func (h *HttpClient) WithToken(token string) *HttpClient {
	h.token = token
	return h
}

func (h *HttpClient) WithBasicAuth(username, password string) *HttpClient {
	h.username = username
	h.password = password
	return h
}

func (h *HttpClient) Timeout(d time.Duration) *HttpClient {
	h.timeout = d
	return h
}

func (h *HttpClient) Retry(count int, delay time.Duration) *HttpClient {
	h.retryCount = count
	h.retryDelay = delay
	return h
}

func (h *HttpClient) WithQuery(key, value string) *HttpClient {
	h.query[key] = value
	return h
}

func (h *HttpClient) ContentType(ct string) *HttpClient {
	h.contentType = ct
	return h
}

func (h *HttpClient) AsJSON() *HttpClient {
	h.contentType = "application/json"
	return h
}

func (h *HttpClient) Get(url string) (*Response, error) {
	return h.doRequest("GET", h.resolveURL(url), nil)
}

func (h *HttpClient) Post(url string, body ...interface{}) (*Response, error) {
	return h.doRequest("POST", h.resolveURL(url), h.encodeBody(body))
}

func (h *HttpClient) Put(url string, body ...interface{}) (*Response, error) {
	return h.doRequest("PUT", h.resolveURL(url), h.encodeBody(body))
}

func (h *HttpClient) Patch(url string, body ...interface{}) (*Response, error) {
	return h.doRequest("PATCH", h.resolveURL(url), h.encodeBody(body))
}

func (h *HttpClient) Delete(url string) (*Response, error) {
	return h.doRequest("DELETE", h.resolveURL(url), nil)
}

func (h *HttpClient) Head(url string) (*Response, error) {
	return h.doRequest("HEAD", h.resolveURL(url), nil)
}

func (h *HttpClient) resolveURL(urlStr string) string {
	if h.baseURL != "" && !strings.HasPrefix(urlStr, "http") {
		return h.baseURL + "/" + strings.TrimLeft(urlStr, "/")
	}
	return urlStr
}

func (h *HttpClient) encodeBody(body []interface{}) io.Reader {
	if len(body) == 0 {
		return nil
	}
	switch v := body[0].(type) {
	case string:
		return strings.NewReader(v)
	case []byte:
		return bytes.NewReader(v)
	case io.Reader:
		return v
	default:
		if h.contentType == "application/json" {
			data, _ := json.Marshal(v)
			return bytes.NewReader(data)
		}
		return strings.NewReader(fmt.Sprintf("%v", v))
	}
}

func (h *HttpClient) applyConfig(req *http.Request) {
	for k, v := range h.headers {
		req.Header.Set(k, v)
	}
	if h.token != "" {
		req.Header.Set("Authorization", "Bearer "+h.token)
	}
	if h.username != "" {
		req.SetBasicAuth(h.username, h.password)
	}
	if h.contentType != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", h.contentType)
	}
	if len(h.query) > 0 {
		q := req.URL.Query()
		for k, v := range h.query {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
}

func (h *HttpClient) doRequest(method, urlStr string, body io.Reader) (*Response, error) {
	var bodyBuf []byte
	if body != nil {
		bodyBuf, _ = io.ReadAll(body)
	}

	client := h.client()
	var lastErr error

	for attempt := 0; attempt <= h.retryCount; attempt++ {
		if attempt > 0 && h.retryDelay > 0 {
			time.Sleep(h.retryDelay)
		}

		var reqBody io.Reader
		if bodyBuf != nil {
			reqBody = bytes.NewReader(bodyBuf)
		}

		req, err := http.NewRequest(method, urlStr, reqBody)
		if err != nil {
			return nil, err
		}
		h.applyConfig(req)

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		bodyBytes, readErr := io.ReadAll(resp.Body)
		resp.Body.Close()
		if readErr != nil {
			lastErr = readErr
			continue
		}

		headers := make(map[string]string)
		for k := range resp.Header {
			headers[k] = resp.Header.Get(k)
		}

		return &Response{
			statusCode: resp.StatusCode,
			body:       string(bodyBytes),
			headers:    headers,
		}, nil
	}

	return nil, lastErr
}

func (h *HttpClient) client() *http.Client {
	if h.timeout > 0 {
		return &http.Client{Timeout: h.timeout}
	}
	return &http.Client{}
}

func (r *Response) StatusCode() int {
	return r.statusCode
}

func (r *Response) Body() string {
	return r.body
}

func (r *Response) JSON(v interface{}) error {
	return json.Unmarshal([]byte(r.body), v)
}

func (r *Response) Headers() map[string]string {
	return r.headers
}

func (r *Response) IsSuccessful() bool {
	return r.statusCode >= 200 && r.statusCode < 300
}

func (r *Response) IsFailed() bool {
	return !r.IsSuccessful()
}

func (r *Response) IsServerError() bool {
	return r.statusCode >= 500
}

func (r *Response) IsClientError() bool {
	return r.statusCode >= 400 && r.statusCode < 500
}
