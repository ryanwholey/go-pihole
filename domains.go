package pihole

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type Domains interface {
	// List returns a list of domains
	List(ctx context.Context, opts *DomainListOptions) (DomainList, error)

	// Get returns a domain by name
	Get(ctx context.Context, domain string, listType string) (*Domain, error)

	// Create adds a domain to the allow or deny lists
	Create(ctx context.Context, opts DomainCreateOptions) (*Domain, error)

	// Delete removes a domain from the allow or deny list
	Delete(ctx context.Context, domain *Domain) (*Domain, error)
}

var (
	DomainTypeWhiteList      = 0
	DomainTypeBlackList      = 1
	DomainTypeWhiteListRegex = 2
	DomainTypeBlackListRegex = 3
)

var AllowedDomainListTypes = []string{
	"black", "white", "regex_black", "regex_white",
}

type domains struct {
	client *Client
}

type DomainList []*Domain

type Domain struct {
	ID      int
	Type    string
	Domain  string
	Enabled bool
	Comment string
}

type domainResponse struct {
	ID      int    `json:"id"`
	Type    int    `json:"type"`
	Domain  string `json:"domain"`
	Enabled int    `json:"enabled"`
	Comment string `json:"comment"`
}

func intToDomainType(i int) (string, error) {
	switch i {
	case 0:
		return "white", nil
	case 1:
		return "black", nil
	case 2:
		return "regex_white", nil
	case 3:
		return "regex_black", nil
	default:
		return "", fmt.Errorf("type %d does not match a known type", i)
	}
}

type domainResponseList []*domainResponse

type domainResponseData struct {
	Data domainResponseList `json:"data"`
}

func (list domainResponseList) ToDomainList() (DomainList, error) {
	l := DomainList{}

	for _, item := range list {
		d, err := item.toDomain()
		if err != nil {
			return nil, err
		}

		l = append(l, d)
	}

	return l, nil
}

func (res domainResponse) toDomain() (*Domain, error) {
	t, err := intToDomainType(res.Type)
	if err != nil {
		return nil, err
	}

	return &Domain{
		ID:      res.ID,
		Type:    t,
		Domain:  res.Domain,
		Enabled: res.Enabled == 1,
		Comment: res.Comment,
	}, nil
}

type DomainListOptions struct {
	ListTypes []string
}

func (d domains) List(ctx context.Context, opts *DomainListOptions) (DomainList, error) {
	var listTypes []string

	if opts != nil {
		listTypes = opts.ListTypes
	} else {
		listTypes = AllowedDomainListTypes
	}

	var responses domainResponseList
	for _, t := range listTypes {
		req, err := d.client.Request(ctx, url.Values{
			"list": []string{t},
		})
		if err != nil {
			return nil, fmt.Errorf("failed fetching list of type %s: %w", t, err)
		}

		res, err := d.client.http.Do(req)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		var data domainResponseData
		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("failed to parse domain list response body: %w", err)
		}

		responses = append(responses, data.Data...)
	}

	return responses.ToDomainList()
}

var (
	ErrDomainNotFound = errors.New("domain not found")
)

// Get returns a domain entry of the matching name and list type
func (d domains) Get(ctx context.Context, domain string, listType string) (*Domain, error) {
	list, err := d.List(ctx, &DomainListOptions{
		ListTypes: []string{listType},
	})
	if err != nil {
		return nil, err
	}

	for _, d := range list {
		if d.Domain == strings.ToLower(domain) {
			return d, nil
		}
	}

	return nil, fmt.Errorf("%w: domain %s of type %s", ErrDomainNotFound, domain, listType)
}

type DomainCreateOptions struct {
	Domain  string
	Comment string
	Enabled bool
	Type    string
}

func (d domains) Create(ctx context.Context, opts DomainCreateOptions) (*Domain, error) {
	return nil, nil
}

func (d domains) Delete(ctx context.Context, domain *Domain) (*Domain, error) {
	return nil, nil
}
