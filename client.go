package pihole

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

type Config struct {
	BaseURL    string
	APIToken   string
	HttpClient *http.Client
	Headers    http.Header
}

type Client struct {
	baseURL    string
	apiToken   string
	headers    http.Header
	http       *http.Client
	LocalDNS   LocalDNS
	LocalCNAME LocalCNAME
	AdBlocker  AdBlocker
	Version    Version
}

// New returns a new Pi-hole client
func New(config Config) *Client {
	baseURL := strings.TrimSuffix(config.BaseURL, "/")

	baseURL = fmt.Sprintf("%s/admin/api.php", baseURL)

	httpClient := retryablehttp.NewClient().StandardClient()
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

	client.LocalDNS = &localDNS{client: client}
	client.LocalCNAME = &localCNAME{client: client}
	client.AdBlocker = &adBlocker{client: client}
	client.Version = &version{client: client}

	return client
}

var ErrClientValidation = errors.New("invalid client configuration")

func (c Client) Validate() error {
	if c.apiToken == "" {
		return fmt.Errorf("%w: apiToken is empty", ErrClientValidation)
	}
	if c.baseURL == "/admin/api.php" {
		return fmt.Errorf("%w: baseURL is empty", ErrClientValidation)
	}

	return nil
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
