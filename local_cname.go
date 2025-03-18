package pihole

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type LocalCNAME interface {
	// List all CNAME records.
	List(ctx context.Context) (CNAMERecordList, error)

	// Create a CNAME record.
	Create(ctx context.Context, domain string, target string) (*CNAMERecord, error)

	// Get a CNAME record by its domain.
	Get(ctx context.Context, domain string) (*CNAMERecord, error)

	// Delete a CNAME record by its domain.
	Delete(ctx context.Context, domain string) error
}

var (
	ErrorLocalCNAMENotFound = errors.New("local CNAME record not found")
)

type localCNAME struct {
	client *Client
}

type CNAMERecord struct {
	Domain string
	Target string
	TTL    int
}

type cnameRecordResponse struct{}

type cnameRecordListResponse struct {
	Config cnameRecordConfigListResponse `json:"config"`
}

type cnameRecordConfigListResponse struct {
	DNS cnameRecordDNSListResponse `json:"dns"`
}

type cnameRecordDNSListResponse struct {
	CNAMERecords []string `json:"cnameRecords"`
}

func (res cnameRecordListResponse) toCNAMERecordList() CNAMERecordList {
	list := make(CNAMERecordList, len(res.Config.DNS.CNAMERecords))

	for i, record := range res.Config.DNS.CNAMERecords {
		entry := strings.Split(record, ",")

		r := CNAMERecord{
			Domain: entry[0],
			Target: entry[1],
		}

		if len(entry) == 3 {
			// TODO: Handle TTL parse error
			ttl, _ := strconv.Atoi(entry[2])
			r.TTL = ttl
		}

		list[i] = r
	}

	return list
}

type CNAMERecordList []CNAMERecord

// List returns all CNAME records
func (cname localCNAME) List(ctx context.Context) (CNAMERecordList, error) {
	res, err := cname.client.Get(ctx, "/api/config/dns/cnameRecords")
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var resList *cnameRecordListResponse
	if err := json.NewDecoder(res.Body).Decode(&resList); err != nil {
		return nil, fmt.Errorf("failed to parse custom CNAME list body: %w", err)
	}

	return resList.toCNAMERecordList(), nil
}

// Create creates a CNAME record
func (cname localCNAME) Create(ctx context.Context, domain string, target string) (*CNAMERecord, error) {
	// TODO: Support TTL
	// bar.com%2Cbaz.com%2C200
	value := fmt.Sprintf("%s%%2C%s", domain, target)
	res, err := cname.client.Put(ctx, fmt.Sprintf("/api/config/dns/cnameRecords/%s", value), nil)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("received unexpected status code %d %s", res.StatusCode, string(b))
	}

	var dnsRes *cnameRecordResponse
	if err := json.NewDecoder(res.Body).Decode(&dnsRes); err != nil {
		return nil, fmt.Errorf("failed to parse custom CNAME response body: %w", err)
	}

	return cname.Get(ctx, domain)
}

// Get returns a CNAME record by the passed domain
func (cname localCNAME) Get(ctx context.Context, domain string) (*CNAMERecord, error) {
	list, err := cname.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch custom CNAME records: %w", err)
	}

	for _, record := range list {
		if strings.EqualFold(record.Domain, domain) {
			return &record, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrorLocalCNAMENotFound, domain)
}

// Delete removes a CNAME record by domain
func (cname localCNAME) Delete(ctx context.Context, domain string) error {
	record, err := cname.Get(ctx, domain)
	if err != nil {
		if errors.Is(err, ErrorLocalCNAMENotFound) {
			return nil
		}
		return fmt.Errorf("failed looking up CNAME record %s for deletion: %w", domain, err)
	}

	value := fmt.Sprintf("%s%%2C%s", record.Domain, record.Target)
	if record.TTL != 0 {
		value = fmt.Sprintf("%s%%2C%d", value, record.TTL)
	}

	res, err := cname.client.Delete(ctx, fmt.Sprintf("/api/config/dns/cnameRecords/%s", value))
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
