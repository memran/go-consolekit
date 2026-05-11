# HTTP Client

Fluent HTTP client with retries, authentication, and JSON support.

## Constructor

```go
func Http() *HttpClient
```

## Configuration

```go
func (h *HttpClient) BaseURL(url string) *HttpClient
func (h *HttpClient) WithHeader(key, value string) *HttpClient
func (h *HttpClient) WithHeaders(headers map[string]string) *HttpClient
func (h *HttpClient) WithToken(token string) *HttpClient
func (h *HttpClient) WithBasicAuth(username, password string) *HttpClient
func (h *HttpClient) Timeout(d time.Duration) *HttpClient
func (h *HttpClient) Retry(count int, delay time.Duration) *HttpClient
func (h *HttpClient) WithQuery(key, value string) *HttpClient
func (h *HttpClient) ContentType(ct string) *HttpClient
func (h *HttpClient) AsJSON() *HttpClient
```

- `AsJSON()` sets `Content-Type: application/json`
- When `BaseURL` is set, relative URLs are resolved against it
- `Retry` retries on network errors up to `count` times with `delay` between attempts

## Request Methods

```go
func (h *HttpClient) Get(url string) (*Response, error)
func (h *HttpClient) Post(url string, body ...interface{}) (*Response, error)
func (h *HttpClient) Put(url string, body ...interface{}) (*Response, error)
func (h *HttpClient) Patch(url string, body ...interface{}) (*Response, error)
func (h *HttpClient) Delete(url string) (*Response, error)
func (h *HttpClient) Head(url string) (*Response, error)
```

Body can be `string`, `[]byte`, `io.Reader`, or any value (JSON-encoded when `AsJSON` is set).

## Response

```go
type Response struct{}

func (r *Response) StatusCode() int
func (r *Response) Body() string
func (r *Response) JSON(v interface{}) error
func (r *Response) Headers() map[string]string
func (r *Response) IsSuccessful() bool
func (r *Response) IsFailed() bool
func (r *Response) IsServerError() bool
func (r *Response) IsClientError() bool
```

## Examples

### GET

```go
resp, err := console.Http().
    BaseURL("https://api.example.com").
    WithHeader("Accept", "application/json").
    WithQuery("page", "1").
    Get("/users")

if resp.IsSuccessful() {
    var users []User
    resp.JSON(&users)
    fmt.Println(resp.Body())
}
```

### POST with JSON

```go
resp, err := console.Http().
    BaseURL("https://api.example.com").
    WithToken("Bearer eyJhbGci...").
    AsJSON().
    Timeout(10 * time.Second).
    Post("/users", map[string]any{
        "name":  "Alice",
        "email": "alice@example.com",
    })

if resp.IsSuccessful() {
    var result map[string]any
    resp.JSON(&result)
}
```

### POST with form data

```go
resp, err := console.Http().
    Post("/submit", "name=Alice&role=admin")
```

### With retries

```go
resp, err := console.Http().
    BaseURL("https://api.example.com").
    Retry(3, 1*time.Second).
    Timeout(5*time.Second).
    Get("/data")
```

### Basic Auth

```go
resp, err := console.Http().
    BaseURL("https://api.example.com").
    WithBasicAuth("admin", "password123").
    Get("/admin/users")
```

### Headers

```go
resp, err := console.Http().
    BaseURL("https://api.example.com").
    ContentType("application/xml").
    WithHeaders(map[string]string{
        "X-Custom": "value",
        "X-API-Key": "abc123",
    }).
    Post("/data", "<xml>...</xml>")

// Response inspection
fmt.Println(resp.StatusCode())  // 200
fmt.Println(resp.Headers()["Content-Type"])
resp.IsSuccessful()             // true
resp.IsFailed()                 // false
resp.IsServerError()            // false
resp.IsClientError()            // false
```

### DELETE

```go
resp, err := console.Http().
    BaseURL("https://api.example.com").
    WithToken("token123").
    Delete("/users/42")
```
