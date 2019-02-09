package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var (
	pathDomainRecords       = "%s/v1/domains/%s/records"
	pathDomainRecordsByType = "%s/v1/domains/%s/records/%s"
	pathDomains             = "%s/v1/domains/%s"
)

// GetDomains fetches the details for the provided domain
func (c *Client) GetDomains(customerID string) ([]Domain, error) {
	domainURL := fmt.Sprintf(pathDomains, c.baseURL, "")
	req, err := http.NewRequest(http.MethodGet, domainURL, nil)

	if err != nil {
		return nil, err
	}

	var d []Domain
	if err := c.execute(customerID, req, &d); err != nil {
		return nil, err
	}

	return d, nil
}

// GetDomain fetches the details for the provided domain
func (c *Client) GetDomain(customerID, domain string) (*Domain, error) {
	domainURL := fmt.Sprintf(pathDomains, c.baseURL, domain)
	req, err := http.NewRequest(http.MethodGet, domainURL, nil)

	if err != nil {
		return nil, err
	}

	d := new(Domain)
	if err := c.execute(customerID, req, &d); err != nil {
		return nil, err
	}

	return d, nil
}

// GetDomainRecords fetches all of the existing records for the provided domain
func (c *Client) GetDomainRecords(customerID, domain string) ([]*DomainRecord, error) {
	domainURL := fmt.Sprintf(pathDomainRecords, c.baseURL, domain)
	req, err := http.NewRequest(http.MethodGet, domainURL, nil)

	if err != nil {
		return nil, err
	}

	records := make([]*DomainRecord, 0)
	if err := c.execute(customerID, req, &records); err != nil {
		return nil, err
	}

	return records, nil
}

// UpdateDomainRecords replaces all of the existing records for the provided domain
func (c *Client) UpdateDomainRecords(customerID, domain string, records []*DomainRecord) error {
	for _, t := range supportedTypes {
		typeRecords := c.domainRecordsOfType(t, records)
		if IsDisallowed(t, typeRecords) {
			continue
		}

		msg, err := json.Marshal(typeRecords)
		if err != nil {
			return err
		}

		buffer := bytes.NewBuffer(msg)
		domainURL := fmt.Sprintf(pathDomainRecordsByType, c.baseURL, domain, t)
		log.Println(domainURL)
		log.Println(buffer)

		req, err := http.NewRequest(http.MethodPut, domainURL, buffer)
		if err != nil {
			return err
		}

		if err := c.execute(customerID, req, nil); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) domainRecordsOfType(t string, records []*DomainRecord) []*DomainRecord {
	typeRecords := make([]*DomainRecord, 0)

	for _, record := range records {
		if strings.EqualFold(record.Type, t) {
			typeRecords = append(typeRecords, record)
		}
	}

	return typeRecords
}
