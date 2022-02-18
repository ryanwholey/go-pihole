package pihole

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type CustomDNSer interface {
	// List all the policies for a given organization
	List(ctx context.Context) (CustomDNSList, error)

	// Create a policy and associate it with an organization.
	Create(ctx context.Context, domain string, IP string) (*CustomDNS, error)

	// Read a policy by its ID.
	Read(ctx context.Context, domain string) (*CustomDNS, error)

	// Update an existing policy.
	Update(ctx context.Context, domain string, IP string) (*CustomDNS, error)

	// Delete a policy by its ID.
	Delete(ctx context.Context, domain string) error
}

// policies implements Policies.
type customDNS struct {
	client *Client
}

type CustomDNS struct {
	IP     string
	Domain string
}

type CustomDNSList []CustomDNS

type CustomDNSListResponse struct {
	Data []CustomDNSResponseObject `json:"data"`
}

type CustomDNSResponse struct {
	Success       string `json:"success"`
	Message       string `json:"message"`
	FTLNotRunning bool   `json:"FTLnotrunning"`
}

type CustomDNSResponseObject []string

func (record CustomDNSResponseObject) ToCustomDNS() CustomDNS {
	return CustomDNS{
		Domain: record[0],
		IP:     record[1],
	}
}

func (res CustomDNSListResponse) ToCustomDNSList() CustomDNSList {
	list := make(CustomDNSList, len(res.Data))

	for i, record := range res.Data {
		list[i] = record.ToCustomDNS()
	}

	return list
}

func (dns customDNS) List(ctx context.Context) (CustomDNSList, error) {
	req, err := dns.client.Request(ctx, url.Values{
		"customdns": []string{"true"},
		"action":    []string{"get"},
	})
	if err != nil {
		return nil, err
	}

	res, err := dns.client.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var resList *CustomDNSListResponse
	if err := json.NewDecoder(res.Body).Decode(&resList); err != nil {
		return nil, fmt.Errorf("failed to parse customDNS list body: %w", err)
	}

	return resList.ToCustomDNSList(), nil
}

func (dns customDNS) Create(ctx context.Context, domain string, IP string) (*CustomDNS, error) {
	req, err := dns.client.Request(ctx, url.Values{
		"customdns": []string{"true"},
		"action":    []string{"add"},
		"ip":        []string{IP},
		"domain":    []string{domain},
	})
	if err != nil {
		return nil, err
	}

	res, err := dns.client.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var dnsRes *CustomDNSResponse
	if err := json.NewDecoder(res.Body).Decode(&dnsRes); err != nil {
		return nil, fmt.Errorf("failed to parse customDNS response body: %w", err)
	}

	record := dnsRes.ToCustomDNS()

	return &record, nil
}

func (dns customDNS) Read(ctx context.Context, domain string) (*CustomDNS, error) {
	return nil, nil
}

func (dns customDNS) Update(ctx context.Context, domain string, IP string) (*CustomDNS, error) {
	return nil, nil
}

func (dns customDNS) Delete(ctx context.Context, domain string) error {
	return nil
}
