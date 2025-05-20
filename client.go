package pihole

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/hashicorp/go-retryablehttp"
)

type Config struct {
	BaseURL    string
	Password   string
	SessionID  string
	HttpClient *http.Client
	Headers    http.Header
}

type Client struct {
	baseURL         string
	password        string
	headers         http.Header
	http            *http.Client
	auth            auth
	publicEndpoints map[string]bool

	sessionLock sync.RWMutex

	LocalDNS   LocalDNS
	LocalCNAME LocalCNAME
	SessionAPI SessionAPI
}

type auth struct {
	sid string
}

const (
	authHeader = "X-FTL-SID"
)

// New returns a new Pi-hole client
func New(config Config) (*Client, error) {
	baseURL := strings.TrimSuffix(config.BaseURL, "/")

	var httpClient *http.Client
	if config.HttpClient != nil {
		httpClient = config.HttpClient
	} else {
		httpClient = retryablehttp.NewClient().StandardClient()
	}

	headers := make(http.Header)
	headers.Add("user-agent", "go-pihole")

	if config.Headers != nil {
		for key, header := range config.Headers {
			headers[key] = header
		}
	}

	client := &Client{
		baseURL:  baseURL,
		http:     httpClient,
		headers:  headers,
		password: config.Password,
		publicEndpoints: map[string]bool{
			"POST /api/auth": true,
		},
	}

	if config.SessionID != "" {
		client.auth.sid = config.SessionID
	}

	client.LocalDNS = &localDNS{client: client}
	client.LocalCNAME = &localCNAME{client: client}
	client.SessionAPI = &sessionAPI{client: client}

	return client, nil
}

var ErrClientValidation = errors.New("invalid client configuration")

func (c *Client) request(ctx context.Context, method string, path string, body interface{}) (*http.Response, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create req with context %s %s: %w", method, path, err)
	}

	if _, ok := c.publicEndpoints[fmt.Sprintf("%s %s", method, path)]; !ok {
		c.sessionLock.RLock()
		SID := c.auth.sid
		c.sessionLock.RUnlock()

		if SID == "" {
			c.sessionLock.Lock()
			defer c.sessionLock.Unlock()

			// recheck client directly to make sure
			if c.auth.sid == "" {
				if _, err := c.SessionAPI.Login(ctx); err != nil {
					return nil, fmt.Errorf("failed to login: %w", err)
				}
			}
		}
		req.Header[authHeader] = []string{c.auth.sid}
	}

	for key, header := range c.headers {
		req.Header[key] = header
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	res, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request %s %s: %w", method, path, err)
	}

	return res, nil
}

func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	return c.request(ctx, "GET", path, nil)
}

func (c *Client) Post(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.request(ctx, "POST", path, body)
}

func (c *Client) Put(ctx context.Context, path string, body interface{}) (*http.Response, error) {
	return c.request(ctx, http.MethodPut, path, body)
}

func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	return c.request(ctx, http.MethodDelete, path, nil)
}

func (c *Client) Request(ctx context.Context, vals url.Values) (*http.Request, error) {
	url := fmt.Sprintf("%s?%s", c.baseURL, vals.Encode())

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, header := range c.headers {
		req.Header[key] = header
	}

	return req, nil
}
