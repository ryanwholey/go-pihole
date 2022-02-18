package pihole

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type CustomDNS interface {
	// List all DNS records.
	List(ctx context.Context) (DNSRecordList, error)

	// Create a DNS record.
	Create(ctx context.Context, domain string, IP string) (*DNSRecord, error)

	// Read a DNS record by its domain.
	Read(ctx context.Context, domain string) (*DNSRecord, error)

	// Update an existing DNS record.
	Update(ctx context.Context, domain string, IP string) (*DNSRecord, error)

	// Delete a DNS record by its domain.
	Delete(ctx context.Context, domain string) error
}

var (
	ErrorCustomDNSNotFound = errors.New("custom dns record not found")
)

type customDNS struct {
	client *Client
}

type DNSRecord struct {
	IP     string
	Domain string
}

type DNSRecordList []DNSRecord

type dnsRecordListResponse struct {
	Data []dnsRecordResponseObject `json:"data"`
}

type dnsRecordResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	FTLNotRunning bool   `json:"FTLnotrunning"`
}

type dnsRecordResponseObject []string

func (record dnsRecordResponseObject) toDNSRecord() DNSRecord {
	return DNSRecord{
		Domain: record[0],
		IP:     record[1],
	}
}

func (res dnsRecordListResponse) toDNSRecordList() DNSRecordList {
	list := make(DNSRecordList, len(res.Data))

	for i, record := range res.Data {
		list[i] = record.toDNSRecord()
	}

	return list
}

// List returns a list of custom DNS records
func (dns customDNS) List(ctx context.Context) (DNSRecordList, error) {
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

	var resList *dnsRecordListResponse
	if err := json.NewDecoder(res.Body).Decode(&resList); err != nil {
		return nil, fmt.Errorf("failed to parse customDNS list body: %w", err)
	}

	return resList.toDNSRecordList(), nil
}

// Create creates a custom DNS record
func (dns customDNS) Create(ctx context.Context, domain string, IP string) (*DNSRecord, error) {
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

	var dnsRes *dnsRecordResponse
	if err := json.NewDecoder(res.Body).Decode(&dnsRes); err != nil {
		return nil, fmt.Errorf("failed to parse customDNS response body: %w", err)
	}

	if !dnsRes.Success {
		return nil, fmt.Errorf("failed to create DNS record %s %s : %s : %w", domain, IP, dnsRes.Message, err)
	}

	return &DNSRecord{
		Domain: domain,
		IP:     IP,
	}, nil
}

// Read returns a custom DNS record by its domain name
func (dns customDNS) Read(ctx context.Context, domain string) (*DNSRecord, error) {
	list, err := dns.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch custom DNS records: %w", err)
	}

	for _, record := range list {
		if record.Domain == domain {
			return &record, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrorCustomDNSNotFound, domain)
}

// Update deletes and recreates a custom DNS record
func (dns customDNS) Update(ctx context.Context, domain string, IP string) (*DNSRecord, error) {
	_, err := dns.Read(ctx, domain)
	if err != nil {
		return nil, err
	}

	if err := dns.Delete(ctx, domain); err != nil {
		return nil, fmt.Errorf("failed to update %s", domain)
	}

	record, err := dns.Create(ctx, domain, IP)
	if err != nil {
		return nil, fmt.Errorf("failed to recreate record during update process: %w", err)
	}

	return record, nil
}

// Delete removes a custom DNS record
func (dns customDNS) Delete(ctx context.Context, domain string) error {
	record, err := dns.Read(ctx, domain)
	if err != nil {
		if errors.Is(err, ErrorCustomDNSNotFound) {
			return nil
		}
		return fmt.Errorf("failed looking up custom DNS record %s for deletion: %w", domain, err)
	}

	req, err := dns.client.Request(ctx, url.Values{
		"customdns": []string{"true"},
		"action":    []string{"delete"},
		"domain":    []string{domain},
		"ip":        []string{record.IP},
	})
	if err != nil {
		return err
	}

	res, err := dns.client.http.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var delRes dnsRecordResponse
	if err := json.NewDecoder(res.Body).Decode(&delRes); err != nil {
		return fmt.Errorf("failed to parse custom DNS deletion response body: %w", err)
	}

	if !delRes.Success {
		return fmt.Errorf("failed to delete custom DNS record %s: %s", domain, delRes.Message)
	}

	return nil
}
