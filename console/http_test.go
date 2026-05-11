package console

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHttpGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Fatalf("expected GET, got %s", r.Method)
		}
		w.Write([]byte(`{"message":"ok"}`))
	}))
	defer ts.Close()

	resp, err := Http().Get(ts.URL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if !resp.IsSuccessful() {
		t.Fatalf("expected success, got %d", resp.StatusCode())
	}
	if resp.Body() != `{"message":"ok"}` {
		t.Fatalf("unexpected body: %s", resp.Body())
	}
}

func TestHttpPostJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("expected JSON content type")
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id":1}`))
	}))
	defer ts.Close()

	resp, err := Http().
		AsJSON().
		Post(ts.URL, map[string]string{"name": "test"})
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
	if resp.StatusCode() != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode())
	}
	var data map[string]interface{}
	if err := resp.JSON(&data); err != nil {
		t.Fatalf("JSON parse failed: %v", err)
	}
	if data["id"] != float64(1) {
		t.Fatalf("expected id 1, got %v", data["id"])
	}
}

func TestHttpPut(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		w.Write([]byte(`{"updated":true}`))
	}))
	defer ts.Close()

	resp, err := Http().Put(ts.URL, `{"name":"updated"}`)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}
	if !resp.IsSuccessful() {
		t.Fatalf("expected success, got %d", resp.StatusCode())
	}
}

func TestHttpPatch(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Fatalf("expected PATCH, got %s", r.Method)
		}
		w.Write([]byte(`{"patched":true}`))
	}))
	defer ts.Close()

	resp, err := Http().Patch(ts.URL, `{"name":"patched"}`)
	if err != nil {
		t.Fatalf("Patch failed: %v", err)
	}
	if !resp.IsSuccessful() {
		t.Fatalf("expected success, got %d", resp.StatusCode())
	}
}

func TestHttpDelete(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	resp, err := Http().Delete(ts.URL)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
	if resp.StatusCode() != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode())
	}
}

func TestHttpHead(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "HEAD" {
			t.Fatalf("expected HEAD, got %s", r.Method)
		}
		w.Header().Set("X-Custom", "value")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	resp, err := Http().Head(ts.URL)
	if err != nil {
		t.Fatalf("Head failed: %v", err)
	}
	if resp.Headers()["X-Custom"] != "value" {
		t.Fatalf("expected X-Custom header 'value', got '%s'", resp.Headers()["X-Custom"])
	}
}

func TestHttpWithHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-Key") != "secret" {
			t.Fatalf("expected X-API-Key header 'secret'")
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Fatalf("expected Accept header")
		}
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	_, err := Http().
		WithHeader("X-API-Key", "secret").
		WithHeader("Accept", "application/json").
		Get(ts.URL)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

func TestHttpWithToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer mytoken" {
			t.Fatalf("expected Bearer token")
		}
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	_, err := Http().WithToken("mytoken").Get(ts.URL)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

func TestHttpWithBasicAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != "admin" || pass != "pass" {
			t.Fatalf("expected basic auth admin:pass")
		}
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	_, err := Http().WithBasicAuth("admin", "pass").Get(ts.URL)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

func TestHttpBaseURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/users" {
			t.Fatalf("expected /api/users, got %s", r.URL.Path)
		}
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	_, err := Http().
		BaseURL(ts.URL).
		Get("/api/users")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

func TestHttpWithQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "2" {
			t.Fatalf("expected page=2, got %s", r.URL.Query().Get("page"))
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Fatalf("expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	_, err := Http().
		WithQuery("page", "2").
		WithQuery("limit", "10").
		Get(ts.URL)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
}

func TestHttpTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	_, err := Http().
		Timeout(10 * time.Millisecond).
		Get(ts.URL)
	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestHttpRetry(t *testing.T) {
	var attempt int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt++
		if attempt == 1 {
			hj, ok := w.(http.Hijacker)
			if !ok {
				t.Fatal("server does not support hijack")
			}
			conn, _, err := hj.Hijack()
			if err != nil {
				t.Fatalf("hijack failed: %v", err)
			}
			conn.Close()
			return
		}
		w.Write([]byte(`{"success":true}`))
	}))
	defer ts.Close()

	resp, err := Http().
		Retry(2, 5*time.Millisecond).
		Get(ts.URL)
	if err != nil {
		t.Fatalf("request failed after retries: %v", err)
	}
	var data map[string]interface{}
	if err := resp.JSON(&data); err != nil {
		t.Fatalf("JSON parse failed: %v", err)
	}
	if data["success"] != true {
		t.Fatal("expected success")
	}
}

func TestHttpResponseStatusChecks(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	resp, _ := Http().Get(ts.URL)
	if resp.IsSuccessful() {
		t.Fatal("expected IsSuccessful=false")
	}
	if !resp.IsFailed() {
		t.Fatal("expected IsFailed=true")
	}
	if !resp.IsClientError() {
		t.Fatal("expected IsClientError=true")
	}
	if resp.IsServerError() {
		t.Fatal("expected IsServerError=false")
	}
}

func TestHttpResponseHeaders(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-RateLimit", "100")
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	resp, _ := Http().Get(ts.URL)
	if resp.Headers()["Content-Type"] != "application/json" {
		t.Fatalf("expected JSON content type, got '%s'", resp.Headers()["Content-Type"])
	}
	if resp.Headers()["X-Ratelimit"] != "100" {
		t.Fatalf("expected rate limit 100, got '%s'", resp.Headers()["X-Ratelimit"])
	}
}

func TestHttpPostStringBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		if string(body) != "raw body" {
			t.Fatalf("unexpected body: %s", string(body))
		}
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	_, err := Http().Post(ts.URL, "raw body")
	if err != nil {
		t.Fatalf("Post failed: %v", err)
	}
}

func TestHttpPutNoBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Fatalf("expected PUT, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	_, err := Http().Put(ts.URL)
	if err != nil {
		t.Fatalf("Put failed: %v", err)
	}
}

func TestHttpGetWithCustomContentType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/xml" {
			t.Fatalf("expected XML content type")
		}
		w.Write([]byte(`<ok/>`))
	}))
	defer ts.Close()

	_, err := Http().
		ContentType("application/xml").
		Get(ts.URL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
}

func TestHttpWithHeadersMap(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-One") != "1" || r.Header.Get("X-Two") != "2" {
			t.Fatalf("unexpected headers")
		}
		w.Write([]byte(`{}`))
	}))
	defer ts.Close()

	_, err := Http().
		WithHeaders(map[string]string{"X-One": "1", "X-Two": "2"}).
		Get(ts.URL)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
}

func TestHttpGetNonExistentHost(t *testing.T) {
	_, err := Http().
		Timeout(100 * time.Millisecond).
		Get("http://192.0.2.1:1")
	if err == nil {
		t.Fatal("expected error for unreachable host")
	}
}
