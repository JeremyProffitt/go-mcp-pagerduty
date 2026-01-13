package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/jeremyproffitt/go-mcp-pagerduty/internal/auth"
)

const (
	DefaultAPIHost = "https://api.pagerduty.com"
	UserAgent      = "go-mcp-pagerduty/0.1.0"
)

// Client is the PagerDuty API client
type Client struct {
	apiKey     string
	apiHost    string
	httpClient *http.Client
	fromEmail  string
}

// Config holds the client configuration
type Config struct {
	APIKey  string
	APIHost string
}

// NewClient creates a new PagerDuty client
func NewClient(cfg Config) *Client {
	apiHost := cfg.APIHost
	if apiHost == "" {
		apiHost = DefaultAPIHost
	}

	return &Client{
		apiKey:  cfg.APIKey,
		apiHost: strings.TrimSuffix(apiHost, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// NewClientFromEnv creates a new client from environment variables
func NewClientFromEnv() (*Client, error) {
	apiKey := os.Getenv("PAGERDUTY_USER_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("PAGERDUTY_USER_API_KEY environment variable is required")
	}

	apiHost := os.Getenv("PAGERDUTY_API_HOST")
	if apiHost == "" {
		apiHost = DefaultAPIHost
	}

	return NewClient(Config{
		APIKey:  apiKey,
		APIHost: apiHost,
	}), nil
}

// SetFromEmail sets the From header for requests (used with user tokens)
func (c *Client) SetFromEmail(email string) {
	c.fromEmail = email
}

// buildURL constructs a full URL with query parameters
func (c *Client) buildURL(path string, params map[string]string) string {
	u := c.apiHost + path

	if len(params) > 0 {
		values := url.Values{}
		for k, v := range params {
			values.Add(k, v)
		}
		u += "?" + values.Encode()
	}

	return u
}

// buildURLWithArrayParams constructs a URL with array query parameters
func (c *Client) buildURLWithArrayParams(path string, params map[string][]string) string {
	u := c.apiHost + path

	if len(params) > 0 {
		values := url.Values{}
		for k, vs := range params {
			for _, v := range vs {
				values.Add(k, v)
			}
		}
		u += "?" + values.Encode()
	}

	return u
}

// getAPIKey returns the API key to use, checking context for override
func (c *Client) getAPIKey(ctx context.Context) string {
	if ctx != nil {
		if token, ok := auth.GetPagerDutyToken(ctx); ok && token != "" {
			return token
		}
	}
	return c.apiKey
}

// doRequest performs an HTTP request with proper headers
func (c *Client) doRequest(method, url string, body interface{}) ([]byte, error) {
	return c.doRequestWithContext(context.Background(), method, url, body)
}

// doRequestWithContext performs an HTTP request with proper headers and context support
func (c *Client) doRequestWithContext(ctx context.Context, method, url string, body interface{}) ([]byte, error) {
	var reqBody io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Token token="+c.getAPIKey(ctx))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Set("User-Agent", UserAgent)

	if c.fromEmail != "" {
		req.Header.Set("From", c.fromEmail)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// Get performs a GET request
func (c *Client) Get(path string, params map[string]string) ([]byte, error) {
	url := c.buildURL(path, params)
	return c.doRequest(http.MethodGet, url, nil)
}

// GetWithArrayParams performs a GET request with array parameters
func (c *Client) GetWithArrayParams(path string, params map[string][]string) ([]byte, error) {
	url := c.buildURLWithArrayParams(path, params)
	return c.doRequest(http.MethodGet, url, nil)
}

// Post performs a POST request
func (c *Client) Post(path string, body interface{}) ([]byte, error) {
	url := c.buildURL(path, nil)
	return c.doRequest(http.MethodPost, url, body)
}

// Put performs a PUT request
func (c *Client) Put(path string, body interface{}) ([]byte, error) {
	url := c.buildURL(path, nil)
	return c.doRequest(http.MethodPut, url, body)
}

// Delete performs a DELETE request
func (c *Client) Delete(path string) ([]byte, error) {
	url := c.buildURL(path, nil)
	return c.doRequest(http.MethodDelete, url, nil)
}

// GetJSON performs a GET request and unmarshals the response
func (c *Client) GetJSON(path string, params map[string]string, v interface{}) error {
	data, err := c.Get(path, params)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// PostJSON performs a POST request and unmarshals the response
func (c *Client) PostJSON(path string, body interface{}, v interface{}) error {
	data, err := c.Post(path, body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// PutJSON performs a PUT request and unmarshals the response
func (c *Client) PutJSON(path string, body interface{}, v interface{}) error {
	data, err := c.Put(path, body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Offset int  `json:"offset"`
	Limit  int  `json:"limit"`
	More   bool `json:"more"`
	Total  int  `json:"total"`
}

// Paginate iterates through all pages of a paginated endpoint
func (c *Client) Paginate(path string, params map[string]string, maxResults int, handler func([]byte) (int, error)) error {
	return c.PaginateWithContext(context.Background(), path, params, maxResults, handler)
}

// PaginateWithContext iterates through all pages of a paginated endpoint with context support
func (c *Client) PaginateWithContext(ctx context.Context, path string, params map[string]string, maxResults int, handler func([]byte) (int, error)) error {
	offset := 0
	limit := 100
	totalFetched := 0

	if params == nil {
		params = make(map[string]string)
	}

	for {
		params["offset"] = fmt.Sprintf("%d", offset)
		params["limit"] = fmt.Sprintf("%d", limit)

		data, err := c.GetWithContext(ctx, path, params)
		if err != nil {
			return err
		}

		count, err := handler(data)
		if err != nil {
			return err
		}

		totalFetched += count

		if maxResults > 0 && totalFetched >= maxResults {
			break
		}

		// Check if there are more results
		var pr PaginatedResponse
		if err := json.Unmarshal(data, &pr); err != nil {
			break
		}

		if !pr.More {
			break
		}

		offset += limit
	}

	return nil
}

// GetWithContext performs a GET request with context support
func (c *Client) GetWithContext(ctx context.Context, path string, params map[string]string) ([]byte, error) {
	url := c.buildURL(path, params)
	return c.doRequestWithContext(ctx, http.MethodGet, url, nil)
}

// GetWithArrayParamsContext performs a GET request with array parameters and context support
func (c *Client) GetWithArrayParamsContext(ctx context.Context, path string, params map[string][]string) ([]byte, error) {
	url := c.buildURLWithArrayParams(path, params)
	return c.doRequestWithContext(ctx, http.MethodGet, url, nil)
}

// PostWithContext performs a POST request with context support
func (c *Client) PostWithContext(ctx context.Context, path string, body interface{}) ([]byte, error) {
	url := c.buildURL(path, nil)
	return c.doRequestWithContext(ctx, http.MethodPost, url, body)
}

// PutWithContext performs a PUT request with context support
func (c *Client) PutWithContext(ctx context.Context, path string, body interface{}) ([]byte, error) {
	url := c.buildURL(path, nil)
	return c.doRequestWithContext(ctx, http.MethodPut, url, body)
}

// DeleteWithContext performs a DELETE request with context support
func (c *Client) DeleteWithContext(ctx context.Context, path string) ([]byte, error) {
	url := c.buildURL(path, nil)
	return c.doRequestWithContext(ctx, http.MethodDelete, url, nil)
}

// GetJSONWithContext performs a GET request and unmarshals the response with context support
func (c *Client) GetJSONWithContext(ctx context.Context, path string, params map[string]string, v interface{}) error {
	data, err := c.GetWithContext(ctx, path, params)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// PostJSONWithContext performs a POST request and unmarshals the response with context support
func (c *Client) PostJSONWithContext(ctx context.Context, path string, body interface{}, v interface{}) error {
	data, err := c.PostWithContext(ctx, path, body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// PutJSONWithContext performs a PUT request and unmarshals the response with context support
func (c *Client) PutJSONWithContext(ctx context.Context, path string, body interface{}, v interface{}) error {
	data, err := c.PutWithContext(ctx, path, body)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}
