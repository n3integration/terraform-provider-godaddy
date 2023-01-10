package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

const (
	defaultLimit = 500

	pathDomainRecords       = "%s/v1/domains/%s/records?limit=%d&offset=%d"
	pathDomainRecordsAdd    = "%s/v1/domains/%s/records"
	pathDomainRecordsUpdate = "%s/v1/domains/%s/records/%s/%s"
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

// GetDomainRecords fetches all existing records for the provided domain
func (c *Client) GetDomainRecords(customerID, domain string) ([]*DomainRecord, error) {
	offset := 1
	records := make([]*DomainRecord, 0)
	for {
		page := make([]*DomainRecord, 0)
		domainURL := fmt.Sprintf(pathDomainRecords, c.baseURL, domain, defaultLimit, offset)
		req, err := http.NewRequest(http.MethodGet, domainURL, nil)

		if err != nil {
			return nil, err
		}

		if err := c.execute(customerID, req, &page); err != nil {
			return nil, err
		}
		if len(page) == 0 {
			break
		}
		offset += 1
		records = append(records, page...)
	}

	return records, nil
}

// AddDomainRecords adds records without affecting existing ones on the provided domain
func (c *Client) AddDomainRecords(customerID, domain string, records []*DomainRecord) error {
	for t := range supportedTypes {
		typeRecords := c.domainRecordsOfType(t, records)
		if IsDisallowed(t, typeRecords) {
			continue
		}

		msg, err := json.Marshal(typeRecords)
		if err != nil {
			return err
		}

		buffer := bytes.NewBuffer(msg)
		domainURL := fmt.Sprintf(pathDomainRecordsAdd, c.baseURL, domain)
		log.Println(domainURL)
		log.Println(buffer)

		// set method to patch to only add records
		// for more info check: https://developer.godaddy.com/doc/endpoint/domains#/v1/recordAdd
		req, err := http.NewRequest(http.MethodPatch, domainURL, buffer)
		if err != nil {
			return err
		}

		if err := c.execute(customerID, req, nil); err != nil {
			return err
		}
	}

	return nil
}

// ReplaceDomainRecords overwrites all existing records with the ones provided
func (c *Client) ReplaceDomainRecords(customerID, domain string, records []*DomainRecord) error {
	for t := range supportedTypes {
		typeRecords := c.domainRecordsOfType(t, records)
		if IsDisallowed(t, typeRecords) {
			continue
		}

		msg, err := json.Marshal(typeRecords)
		if err != nil {
			return err
		}

		domainURL := fmt.Sprintf(pathDomainRecordsByType, c.baseURL, domain, t)
		buffer := bytes.NewBuffer(msg)

		log.Println(domainURL)
		log.Println(buffer)

		// set method to put to replace all existing records
		// for more info check: https://developer.godaddy.com/doc/endpoint/domains#/v1/recordReplaceType
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

// AddDomainRecords adds records without affecting existing ones on the provided domain
func (c *Client) UpdateDomainRecords(customerID, domain string, records []*DomainRecord) error {
	for _, rec := range records {
		// typeRecords := c.domainRecordsOfType(t, records)
		t := rec.Type
		// if IsDisallowed(t, typeRecords) {
		// 	continue
		// }

		msg, err := json.Marshal([]*DomainRecord{rec})
		if err != nil {
			return err
		}

		buffer := bytes.NewBuffer(msg)
		domainURL := fmt.Sprintf(pathDomainRecordsUpdate, c.baseURL, domain, t, rec.Name)
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
