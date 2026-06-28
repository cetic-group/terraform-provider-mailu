package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultTimeout = 30 * time.Second

type Client struct {
	endpoint   *url.URL
	token      string
	httpClient *http.Client
}

func New(endpoint string, token string) (*Client, error) {
	if strings.TrimSpace(endpoint) == "" {
		return nil, fmt.Errorf("endpoint is required")
	}

	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse endpoint: %w", err)
	}

	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("endpoint must be an absolute URL")
	}

	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("token is required")
	}

	return &Client{
		endpoint: parsed,
		token:    token,
		httpClient: &http.Client{
			Timeout: defaultTimeout,
		},
	}, nil
}

func (c *Client) Get(ctx context.Context, path string, out any) error {
	return c.do(ctx, http.MethodGet, path, nil, out)
}

func (c *Client) Post(ctx context.Context, path string, in any, out any) error {
	return c.do(ctx, http.MethodPost, path, in, out)
}

func (c *Client) Put(ctx context.Context, path string, in any, out any) error {
	return c.do(ctx, http.MethodPut, path, in, out)
}

func (c *Client) Delete(ctx context.Context, path string) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) do(ctx context.Context, method string, path string, in any, out any) error {
	body, err := encodeBody(in)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, method, c.resolve(path), body)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	if in != nil {
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
		return fmt.Errorf("mailu api returned %s: %s", resp.Status, strings.TrimSpace(string(payload)))
	}

	if out == nil || len(payload) == 0 {
		return nil
	}

	if err := json.Unmarshal(payload, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	return nil
}

func encodeBody(in any) (io.Reader, error) {
	if in == nil {
		return nil, nil
	}

	payload, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("encode request: %w", err)
	}

	return bytes.NewReader(payload), nil
}

func (c *Client) resolve(path string) string {
	base := *c.endpoint
	base.Path = strings.TrimRight(base.Path, "/") + "/" + strings.TrimLeft(path, "/")
	return base.String()
}
