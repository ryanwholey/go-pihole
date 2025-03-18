package pihole

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type LocalDNS interface {
	// List all DNS records.
	List(ctx context.Context) (DNSRecordList, error)

	// Create a DNS record.
	Create(ctx context.Context, domain string, IP string) (*DNSRecord, error)

	// Get a DNS record by its domain.
	Get(ctx context.Context, domain string) (*DNSRecord, error)

	// Delete a DNS record by its domain.
	Delete(ctx context.Context, domain string) error
}

var (
	ErrorLocalDNSNotFound = errors.New("local dns record not found")
)

type localDNS struct {
	client *Client
}

type DNSRecord struct {
	IP     string
	Domain string
}

type DNSRecordList []DNSRecord

type dnsRecordListResponse struct {
	Config dnsRecordConfigListResponse `json:"config"`
}

type dnsRecordConfigListResponse struct {
	DNS dnsRecordDNSListResponse `json:"dns"`
}

type dnsRecordDNSListResponse struct {
	Hosts []string `json:"hosts"`
}

type dnsRecordResponse struct{}

func (res dnsRecordListResponse) toDNSRecordList() DNSRecordList {
	list := make(DNSRecordList, len(res.Config.DNS.Hosts))

	for i, record := range res.Config.DNS.Hosts {
		entry := strings.Split(record, " ")

		list[i] = DNSRecord{IP: entry[0], Domain: entry[1]}
	}

	return list
}

// List returns a list of custom DNS records
func (dns localDNS) List(ctx context.Context) (DNSRecordList, error) {
	res, err := dns.client.Get(ctx, "/api/config/dns/hosts")
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
func (dns localDNS) Create(ctx context.Context, domain string, IP string) (*DNSRecord, error) {
	value := fmt.Sprintf("%s%%20%s", IP, domain)

	res, err := dns.client.Put(ctx, fmt.Sprintf("/api/config/dns/hosts/%s", value), nil)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("received unexpected status code %d %s", res.StatusCode, string(b))
	}

	var dnsRes *dnsRecordResponse
	if err := json.NewDecoder(res.Body).Decode(&dnsRes); err != nil {
		return nil, fmt.Errorf("failed to parse customDNS response body: %w", err)
	}

	// if !dnsRes.Success {
	// 	return nil, fmt.Errorf("failed to create DNS record %s %s : %s : %w", domain, IP, dnsRes.Message, err)
	// }

	return dns.Get(ctx, domain)
}

// Get returns a custom DNS record by its domain name
func (dns localDNS) Get(ctx context.Context, domain string) (*DNSRecord, error) {
	records, err := dns.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch custom DNS records: %w", err)
	}

	for _, record := range records {
		if strings.EqualFold(record.Domain, domain) {
			return &record, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrorLocalDNSNotFound, domain)
}

// Delete removes a custom DNS record
func (dns localDNS) Delete(ctx context.Context, domain string) error {
	record, err := dns.Get(ctx, domain)
	if err != nil {
		if errors.Is(err, ErrorLocalDNSNotFound) {
			return nil
		}

		return fmt.Errorf("failed looking up custom DNS record %s for deletion: %w", domain, err)
	}

	value := fmt.Sprintf("%s%%20%s", record.IP, record.Domain)

	res, err := dns.client.Delete(ctx, fmt.Sprintf("/api/config/dns/hosts/%s", value))
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(res.Body)
		return fmt.Errorf("received unexpected status code: %d %s", res.StatusCode, string(b))
	}

	return nil
}
