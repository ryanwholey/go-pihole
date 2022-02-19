package pihole

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

type LocalCNAME interface {
	// List all CNAME records.
	List(ctx context.Context) (CNAMERecordList, error)

	// Create a CNAME record.
	Create(ctx context.Context, domain string, target string) (*CNAMERecord, error)

	// Get a CNAME record by its domain.
	Get(ctx context.Context, domain string) (*CNAMERecord, error)

	// Update an existing CNAME record.
	Update(ctx context.Context, domain string, IP string) (*CNAMERecord, error)

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
}

type cnameRecordResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	FTLNotRunning bool   `json:"FTLnotrunning"`
}

type cnameRecordListResponse struct {
	Data []cnameRecordResponseObject `json:"data"`
}

func (res cnameRecordListResponse) toCNAMERecordList() CNAMERecordList {
	list := make(CNAMERecordList, len(res.Data))

	for i, record := range res.Data {
		list[i] = record.toCNAMERecord()
	}

	return list
}

type cnameRecordResponseObject []string

func (record cnameRecordResponseObject) toCNAMERecord() CNAMERecord {
	return CNAMERecord{
		Domain: record[0],
		Target: record[1],
	}
}

type CNAMERecordList []CNAMERecord

// List returns all CNAME records
func (cname localCNAME) List(ctx context.Context) (CNAMERecordList, error) {
	req, err := cname.client.Request(ctx, url.Values{
		"customcname": []string{"true"},
		"action":      []string{"get"},
	})
	if err != nil {
		return nil, err
	}

	res, err := cname.client.http.Do(req)
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
	req, err := cname.client.Request(ctx, url.Values{
		"customcname": []string{"true"},
		"action":      []string{"add"},
		"domain":      []string{domain},
		"target":      []string{target},
	})
	if err != nil {
		return nil, err
	}

	res, err := cname.client.http.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var dnsRes *cnameRecordResponse
	if err := json.NewDecoder(res.Body).Decode(&dnsRes); err != nil {
		return nil, fmt.Errorf("failed to parse custom CNAME response body: %w", err)
	}

	if !dnsRes.Success {
		return nil, fmt.Errorf("failed to create CNAME record %s %s : %s : %w", domain, target, dnsRes.Message, err)
	}

	return &CNAMERecord{
		Domain: domain,
		Target: target,
	}, nil
}

// Get returns a CNAME record by the passed domain
func (cname localCNAME) Get(ctx context.Context, domain string) (*CNAMERecord, error) {
	list, err := cname.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch custom CNAME records: %w", err)
	}

	for _, record := range list {
		if record.Domain == domain {
			return &record, nil
		}
	}

	return nil, fmt.Errorf("%w: %s", ErrorLocalCNAMENotFound, domain)
}

// Update deletes and recreates a CNAME record
func (cname localCNAME) Update(ctx context.Context, domain string, target string) (*CNAMERecord, error) {
	_, err := cname.Get(ctx, domain)
	if err != nil {
		return nil, err
	}

	if err := cname.Delete(ctx, domain); err != nil {
		return nil, fmt.Errorf("failed to update %s", domain)
	}

	record, err := cname.Create(ctx, domain, target)
	if err != nil {
		return nil, fmt.Errorf("failed to recreate record during update process: %w", err)
	}

	return record, nil
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

	req, err := cname.client.Request(ctx, url.Values{
		"customcname": []string{"true"},
		"action":      []string{"delete"},
		"domain":      []string{domain},
		"target":      []string{record.Target},
	})
	if err != nil {
		return err
	}

	res, err := cname.client.http.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	var delRes cnameRecordResponse
	if err := json.NewDecoder(res.Body).Decode(&delRes); err != nil {
		return fmt.Errorf("failed to parse CNAME deletion response body: %w", err)
	}

	if !delRes.Success {
		return fmt.Errorf("failed to delete CNAME record %s: %s", domain, delRes.Message)
	}

	return nil
}
