package pihole

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Config struct {
	BaseURL    string
	APIToken   string
	HttpClient *http.Client
	Headers    http.Header
}

type Client struct {
	baseURL   string
	apiToken  string
	headers   http.Header
	http      *http.Client
	CustomDNS CustomDNSer
}

// New returns a new Pi-Hole client
func New(config Config) *Client {
	baseURL := config.BaseURL

	if strings.HasSuffix(baseURL, "/") {
		baseURL = strings.TrimSuffix(baseURL, "/")
	}

	baseURL = fmt.Sprintf("%s/admin/api.php", baseURL)

	httpClient := &http.Client{}
	if config.HttpClient != nil {
		httpClient = config.HttpClient
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
		apiToken: config.APIToken,
		http:     httpClient,
		headers:  headers,
	}

	client.CustomDNS = &customDNS{client: client}

	return client
}

// Get returns a new get request
func (c Client) Get(ctx context.Context, vals url.Values) (*http.Request, error) {
	vals.Set("auth", c.apiToken)

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

func (c Client) Request(ctx context.Context, vals url.Values) (*http.Request, error) {
	vals.Set("auth", c.apiToken)

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
