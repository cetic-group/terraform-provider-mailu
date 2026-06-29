package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestClientGetSuccess(t *testing.T) {
	t.Parallel()

	var authHeader string
	var userAgent string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader = r.Header.Get("Authorization")
		userAgent = r.Header.Get("User-Agent")

		if r.URL.Path != "/api/v1/domain" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"name":"example.com"}`))
	}))
	defer server.Close()

	c, err := NewWithConfig(Config{
		Endpoint:   server.URL + "/api/v1",
		Token:      "secret-token",
		UserAgent:  "test-agent",
		MaxRetries: -1,
	})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	var out map[string]string
	if err := c.Get(context.Background(), "/domain", &out); err != nil {
		t.Fatalf("get: %v", err)
	}

	if got, want := out["name"], "example.com"; got != want {
		t.Fatalf("name = %q, want %q", got, want)
	}
	if got, want := authHeader, "Bearer secret-token"; got != want {
		t.Fatalf("authorization = %q, want %q", got, want)
	}
	if got, want := userAgent, "test-agent"; got != want {
		t.Fatalf("user-agent = %q, want %q", got, want)
	}
}

func TestClientAPIError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":400,"message":"bad raw_password=secret"}`))
	}))
	defer server.Close()

	c, err := NewWithConfig(Config{
		Endpoint:   server.URL,
		Token:      "secret-token",
		MaxRetries: -1,
	})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	err = c.Get(context.Background(), "/domain", nil)
	if err == nil {
		t.Fatal("expected error")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error type = %T, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", apiErr.StatusCode, http.StatusBadRequest)
	}
	if strings.Contains(err.Error(), "secret") {
		t.Fatalf("error was not redacted: %s", err)
	}
}

func TestClientResolveURL(t *testing.T) {
	t.Parallel()

	c, err := New("https://mail.example.com/api/v1/", "token")
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	if got, want := c.Resolve("/domain/example.com"), "https://mail.example.com/api/v1/domain/example.com"; got != want {
		t.Fatalf("resolved URL = %q, want %q", got, want)
	}
}

func TestClientRetriesServerErrors(t *testing.T) {
	t.Parallel()

	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts++
		if attempts == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"message":"temporary"}`))
			return
		}

		_, _ = w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	c, err := NewWithConfig(Config{
		Endpoint:   server.URL,
		Token:      "token",
		MaxRetries: 1,
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatalf("new client: %v", err)
	}

	if err := c.Get(context.Background(), "/domain", &map[string]bool{}); err != nil {
		t.Fatalf("get: %v", err)
	}
	if attempts != 2 {
		t.Fatalf("attempts = %d, want 2", attempts)
	}
}
