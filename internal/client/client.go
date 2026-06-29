package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout    = 30 * time.Second
	defaultMaxRetries = 2
	defaultUserAgent  = "terraform-provider-mailu/dev"
)

type Client struct {
	endpoint   *url.URL
	token      string
	httpClient *http.Client
	userAgent  string
	maxRetries int
}

type Config struct {
	Endpoint              string
	Token                 string
	Timeout               time.Duration
	MaxRetries            int
	UserAgent             string
	InsecureSkipTLSVerify bool
}

type APIError struct {
	StatusCode int
	Status     string
	Code       int
	Message    string
	Body       string
	RetryAfter time.Duration
}

func (e *APIError) Error() string {
	message := strings.TrimSpace(e.Message)
	if message == "" {
		message = strings.TrimSpace(e.Body)
	}
	if message == "" {
		message = "empty response body"
	}

	return fmt.Sprintf("mailu api returned %s: %s", e.Status, Redact(message))
}

func New(endpoint string, token string) (*Client, error) {
	return NewWithConfig(Config{
		Endpoint: endpoint,
		Token:    token,
	})
}

func NewWithConfig(config Config) (*Client, error) {
	if strings.TrimSpace(config.Endpoint) == "" {
		return nil, fmt.Errorf("endpoint is required")
	}

	parsed, err := url.Parse(config.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse endpoint: %w", err)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("endpoint must be an absolute URL")
	}

	if strings.TrimSpace(config.Token) == "" {
		return nil, fmt.Errorf("token is required")
	}

	timeout := config.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	maxRetries := config.MaxRetries
	if maxRetries == 0 {
		maxRetries = defaultMaxRetries
	} else if maxRetries < 0 {
		maxRetries = 0
	}

	userAgent := strings.TrimSpace(config.UserAgent)
	if userAgent == "" {
		userAgent = defaultUserAgent
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	if config.InsecureSkipTLSVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec
	}

	return &Client{
		endpoint: parsed,
		token:    config.Token,
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
		userAgent:  userAgent,
		maxRetries: maxRetries,
	}, nil
}

func (c *Client) Get(ctx context.Context, path string, out any) error {
	return c.do(ctx, http.MethodGet, path, nil, out)
}

func (c *Client) Post(ctx context.Context, path string, in any, out any) error {
	return c.do(ctx, http.MethodPost, path, in, out)
}

func (c *Client) Patch(ctx context.Context, path string, in any, out any) error {
	return c.do(ctx, http.MethodPatch, path, in, out)
}

func (c *Client) Delete(ctx context.Context, path string) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) Resolve(path string) string {
	return c.resolve(path)
}

func (c *Client) do(ctx context.Context, method string, path string, in any, out any) error {
	body, err := encodeBody(in)
	if err != nil {
		return err
	}

	attempts := c.maxRetries + 1
	for attempt := 1; attempt <= attempts; attempt++ {
		err = c.doOnce(ctx, method, path, body, in != nil, out)
		if err == nil {
			return nil
		}

		if !isRetryable(method, err) || attempt == attempts {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(retryDelay(err, attempt)):
		}
	}

	return err
}

func (c *Client) doOnce(ctx context.Context, method string, path string, body []byte, hasBody bool, out any) error {
	req, err := http.NewRequestWithContext(ctx, method, c.resolve(path), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("User-Agent", c.userAgent)
	if hasBody {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("mailu api request failed: %w", err)
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return newAPIError(resp, payload)
	}

	if out == nil || len(payload) == 0 {
		return nil
	}

	if err := json.Unmarshal(payload, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func encodeBody(in any) ([]byte, error) {
	if in == nil {
		return nil, nil
	}

	payload, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	return payload, nil
}

func (c *Client) resolve(path string) string {
	base := *c.endpoint
	base.Path = strings.TrimRight(base.Path, "/") + "/" + strings.TrimLeft(path, "/")
	return base.String()
}

func newAPIError(resp *http.Response, payload []byte) error {
	apiErr := &APIError{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Body:       strings.TrimSpace(string(payload)),
		RetryAfter: parseRetryAfter(resp.Header.Get("Retry-After"), time.Now()),
	}

	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(payload, &response); err == nil {
		apiErr.Code = response.Code
		apiErr.Message = response.Message
	}

	return apiErr
}

func isRetryable(method string, err error) bool {
	if method != http.MethodGet && method != http.MethodDelete && method != http.MethodHead {
		return false
	}

	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == http.StatusTooManyRequests || apiErr.StatusCode >= 500
	}

	return false
}

func retryDelay(err error, attempt int) time.Duration {
	var apiErr *APIError
	if errors.As(err, &apiErr) && apiErr.RetryAfter > 0 {
		return apiErr.RetryAfter
	}

	return time.Duration(attempt) * 200 * time.Millisecond
}

func parseRetryAfter(value string, now time.Time) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}

	if seconds, err := time.ParseDuration(value + "s"); err == nil {
		return seconds
	}

	when, err := http.ParseTime(value)
	if err != nil {
		return 0
	}

	delay := time.Until(when)
	if !now.IsZero() {
		delay = when.Sub(now)
	}
	if delay < 0 {
		return 0
	}

	return delay
}
