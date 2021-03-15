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
	pathDomainRecordsByTypeName = "%s/v1/domains/%s/records/%s/%s"
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


func (c *Client) GetExistingDomainRecord(domain string, recordType string, name string,  record *DomainRecord) (*DomainRecord, error) {

	domainRecords, err := c.GetDomainRecords("", domain)
	for _, t := range supportedTypes {
		typeRecords := c.domainRecordsOfType(t, domainRecords)
		if IsDisallowed(t, typeRecords) { continue }
		msg, err := json.Marshal(typeRecords)
		buffer := bytes.NewBuffer(msg)
		log.Println(buffer)
		err = err
	}
	return record, err	
}

// UpdateDomainRecords replaces domain records
func (c *Client) UpdateDomainRecords(customerID, domain string, records []*DomainRecord) error {
	domainRecords, err := c.GetDomainRecords("", domain)
	err = err


	// PROCESS RECORDS WITH UNIQUE RECORD TYPE AND NAME, SUMIT THEM ONE BY ONE
	for _, record := range records {
		
		if (record.Type == "A" || record.Type == "AAAA" || record.Type == "CNAME" ) {
			record.Name = strings.ToLower(record.Name)
			record.Data = strings.ToLower(record.Data)
		}

		log.Println("[Trace] EVALUATING NEW RECORD: - TYPE: " + record.Type + " - NAME: " + record.Name + " - DATA: " + record.Data)

		for _, existingRecord := range domainRecords {

			if IsDisallowed(existingRecord.Type, records) { continue }


			if (existingRecord.Type == record.Type && existingRecord.Name == record.Name && existingRecord.Data != record.Data && record.Name != "@") {
				log.Println("[Trace] ################################")
				log.Println("[Trace] FOUND CHANGED RECORD")
				log.Println("[Trace] Existing: " + existingRecord.Name + " - " + existingRecord.Data)
				log.Println("[Trace] New: " + record.Name + " - " + record.Data)
				
				domainURL := fmt.Sprintf(pathDomainRecordsByTypeName, c.baseURL, domain, record.Type, record.Name)
				var recordArray = make([]*DomainRecord, 0)
				recordArray = append(recordArray, record)	
				msg, err := json.Marshal(recordArray)
				buffer := bytes.NewBuffer(msg)
				err = err

				req, err := http.NewRequest(http.MethodPut, domainURL, buffer)

				log.Println(domainURL)
				log.Println(buffer)

				if err != nil { return err }
				if err := c.execute(customerID, req, nil); err != nil {	return err }
			}		
		}
	}

	// PROCESS RECORDS WITH "@" SUBMIT THEM ON A PER TYPE BASIS IN BATCH
	originNameRecords := c.domainRecordsOfName("@", records)

	for _, supportedType := range supportedTypes {
		originNameTypeRecords := c.domainRecordsOfType(supportedType, originNameRecords)
		if (len(originNameTypeRecords) > 0) {
			msg, err := json.Marshal(originNameTypeRecords)
			if err != nil { return err }
			buffer := bytes.NewBuffer(msg)
	
			log.Println("NOW SUBMITTING @ RECORDS OF TYPE " + supportedType)
	
			domainURL := fmt.Sprintf(pathDomainRecordsByTypeName, c.baseURL, domain, supportedType, "@")
			req, err := http.NewRequest(http.MethodPut, domainURL, buffer)
			log.Println(domainURL)
			log.Println(buffer)
	
			if err != nil { return err }
			if err := c.execute(customerID, req, nil); err != nil {	return err }	
		}
	}
	return nil
}

func (c *Client) domainRecordsOfName(name string, records []*DomainRecord) []*DomainRecord {
	typeRecords := make([]*DomainRecord, 0)

	for _, record := range records {
		if strings.EqualFold(record.Name, name) {
			typeRecords = append(typeRecords, record)
		}
	}

	return typeRecords
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
