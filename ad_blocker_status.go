package pihole

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type AdBlocker interface {
	// Get returns the ad blocker status
	Get(ctx context.Context) (*AdBlockerStatus, error)

	// Update updates the ad blocker status (enabled, disabled)
	Update(ctx context.Context, opts AdBlockerStatusOptions) (*AdBlockerStatus, error)
}

type AdBlockerStatusOptions struct {
	Enabled         bool
	DisabledSeconds int
}

type adBlocker struct {
	client *Client
}

type adBlockerStatusResponse struct {
	Status string `json:"status"`
}

// AdBlockerStatus is an object representing the ad blocker status
type AdBlockerStatus struct {
	Enabled bool
}

func (res adBlockerStatusResponse) toAdBlockerStatus() *AdBlockerStatus {
	return &AdBlockerStatus{
		Enabled: strings.EqualFold(res.Status, "enabled"),
	}
}

// Get returns the ad blocker status
func (ab adBlocker) Get(ctx context.Context) (*AdBlockerStatus, error) {
	req, err := ab.client.Request(ctx, url.Values{
		"status": []string{"true"},
	})
	if err != nil {
		return nil, err
	}

	res, err := ab.client.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var status *adBlockerStatusResponse
	if err := json.NewDecoder(res.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to parse ad blocker status body: %w", err)
	}

	return status.toAdBlockerStatus(), nil
}

// Update changes the ad blocker status state
func (ab adBlocker) Update(ctx context.Context, opts AdBlockerStatusOptions) (*AdBlockerStatus, error) {
	action := "enable"
	val := fmt.Sprint(true)

	if !opts.Enabled {
		action = "disable"
		val = fmt.Sprint(opts.DisabledSeconds)
	}

	req, err := ab.client.Request(ctx, url.Values{
		action: []string{val},
	})

	if err != nil {
		return nil, err
	}

	res, err := ab.client.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var status *adBlockerStatusResponse
	if err := json.NewDecoder(res.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to parse ad blocker status body: %w", err)
	}

	return status.toAdBlockerStatus(), nil
}
