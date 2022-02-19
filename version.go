package pihole

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type Version interface {
	// Get returns the Pi-hole server component versions
	Get(ctx context.Context) (*ComponentVersions, error)
}

type version struct {
	client *Client
}

type ComponentVersions struct {
	CoreUpdate  bool   `json:"core_update,omitempty"`
	WebUpdate   bool   `json:"web_update,omitempty"`
	FTLUpdate   bool   `json:"FTL_update,omitempty"`
	CoreCurrent string `json:"core_current,omitempty"`
	WebCurrent  string `json:"web_current,omitempty"`
	FTLCurrent  string `json:"FTL_current,omitempty"`
	CoreLatest  string `json:"core_latest,omitempty"`
	WebLatest   string `json:"web_latest,omitempty"`
	FTLLatest   string `json:"FTL_latest,omitempty"`
	CoreBranch  string `json:"core_branch,omitempty"`
	WebBranch   string `json:"web_branch,omitempty"`
	FTLBranch   string `json:"FTL_branch,omitempty"`
}

func (v version) Get(ctx context.Context) (*ComponentVersions, error) {
	req, err := v.client.Request(ctx, url.Values{
		"versions": []string{"true"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to format request to fetch versions: %w", err)
	}

	res, err := v.client.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch versions: %w", err)
	}

	defer res.Body.Close()

	var vRes *ComponentVersions
	if err := json.NewDecoder(res.Body).Decode(&vRes); err != nil {
		return nil, fmt.Errorf("failed to parse versions response body: %w", err)
	}

	return vRes, nil
}
