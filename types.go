package main

import (
	"fmt"
	"strings"
)

// RecordType is an enumeration of possible DNS record types
type RecordType int

const (
	A RecordType = iota
	AAAA
	CNAME
	MX
	NS
	SOA
	SRV
	TXT
)

var supportedTypes = []string{
	"A", "AAAA", "CNAME", "NS", "SOA", "TXT",
}

// Domain encapsulates a domain resource
type Domain struct {
	ID     int64  `json:"domainId"`
	Name   string `json:"domain"`
	Status string `json:"status"`
}

// DomainRecord encapsulates a domain record resource
type DomainRecord struct {
	Type string `json:"type,omitempty"`
	Name string `json:"name"`
	Data string `json:"data"`
	//Priority int    `json:"priority"`
	TTL int `json:"ttl"`
	//Service  string `json:"service"`
	//Protocol string `json:"protocol"`
	//Port     int    `json:"port"`
	//Weight   int    `json:"weight"`
}

// NewDomainRecord validates and constructs a DomainRecord, if valid.
func NewDomainRecord(name, t, data string, ttl int) (*DomainRecord, error) {
	name = strings.TrimSpace(name)
	data = strings.TrimSpace(data)
	if len(name) < 1 || len(name) > 255 {
		return nil, fmt.Errorf("name must be between 1..255")
	}
	if len(data) < 1 || len(data) > 255 {
		return nil, fmt.Errorf("data must be between 1..255")
	}
	if ttl < 0 {
		return nil, fmt.Errorf("ttl must be a positive value")
	}
	if !isSupportedType(t) {
		return nil, fmt.Errorf("type must be one of: %s", supportedTypes)
	}
	return &DomainRecord{
		Name: name,
		Type: t,
		Data: data,
		TTL:  ttl,
	}, nil
}

func isSupportedType(recType string) bool {
	for _, t := range supportedTypes {
		if t == recType {
			return true
		}
	}
	return false
}
