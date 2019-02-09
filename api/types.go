package api

import (
	"fmt"
	"strings"
)

// RecordType is an enumeration of possible DNS record types
type RecordType int

// RecordFactory is a factory method for creating new DomainRecords
type RecordFactory func(string) (*DomainRecord, error)

const (
	// A is an address record type
	A RecordType = iota
	// AAAA is an IPv6 address record type
	AAAA
	// CNAME is a Canonical record name (alias) type
	CNAME
	// MX is a mail exchange record type
	MX
	// NS is a name server record type
	NS
	// SOA is a start of authority record type
	SOA
	// SRV is a service locator type
	SRV
	// TXT is a text record type
	TXT
)

func (rt RecordType) String() string {
	switch rt {
	case A:
		return AType
	case AAAA:
		return AAAAType
	case CNAME:
		return CNameType
	case MX:
		return MXType
	case NS:
		return NSType
	case SOA:
		return SOAType
	case TXT:
		return TXTType
	}
	return ""
}

const (
	DefaultTTL      = 3600
	DefaultPriority = 0

	StatusActive    = "ACTIVE"
	StatusCancelled = "CANCELLED"

	Ptr       = "@"
	AType     = "A"
	AAAAType  = "AAAA"
	CNameType = "CNAME"
	MXType    = "MX"
	NSType    = "NS"
	SOAType   = "SOA"
	TXTType   = "TXT"
)

var supportedTypes = []string{
	AType, AAAAType, CNameType, MXType, NSType, SOAType, TXTType,
}

// Domain encapsulates a domain resource
type Domain struct {
	ID     int64  `json:"domainId"`
	Name   string `json:"domain"`
	Status string `json:"status"`
}

// DomainRecord encapsulates a domain record resource
type DomainRecord struct {
	Type     string `json:"type,omitempty"`
	Name     string `json:"name"`
	Data     string `json:"data"`
	Priority int    `json:"priority,omitempty"`
	TTL      int    `json:"ttl"`
	// Service  string `json:"service"`
	// Protocol string `json:"protocol"`
	// Port     int    `json:"port"`
	// Weight   int    `json:"weight"`
}

// NewDomainRecord validates and constructs a DomainRecord, if valid.
func NewDomainRecord(name, t, data string, ttl int, priority int) (*DomainRecord, error) {
	name = strings.TrimSpace(name)
	data = strings.TrimSpace(data)
	if err := ValidateData(t, data); err != nil {
		return nil, err
	}

	parts := strings.Split(name, ".")
	if len(parts) < 1 || len(parts) > 255 {
		return nil, fmt.Errorf("name must be between 1..255 octets")
	}
	for _, part := range parts {
		if len(part) > 63 {
			return nil, fmt.Errorf("invalid domain name. name octets should be less than 63 characters")
		}
	}

	if ttl < 0 {
		return nil, fmt.Errorf("ttl must be a positive value")
	}
	if err := ValidatePriority(priority); err != nil {
		return nil, err
	}
	if !isSupportedType(t) {
		return nil, fmt.Errorf("type must be one of: %s", supportedTypes)
	}
	return &DomainRecord{
		Name:     name,
		Type:     t,
		Data:     data,
		TTL:      ttl,
		Priority: priority,
	}, nil
}

// NewNSRecord constructs a nameserver record from the supplied data
func NewNSRecord(data string) (*DomainRecord, error) {
	return NewDomainRecord(Ptr, NSType, data, DefaultTTL, DefaultPriority)
}

// NewARecord constructs a new address record from the supplied data
func NewARecord(data string) (*DomainRecord, error) {
	return NewDomainRecord(Ptr, AType, data, DefaultTTL, DefaultPriority)
}

// ValidateData performs bounds checking on a data element
func ValidateData(t, data string) error {
	switch t {
	case TXTType:
		if len(data) < 0 || len(data) > 512 {
			return fmt.Errorf("data must be between 0..512 characters in length")
		}
	default:
		if len(data) < 0 || len(data) > 255 {
			return fmt.Errorf("data must be between 0..255 characters in length")
		}
	}
	return nil
}

// ValidatePriority performs bounds checking on priority element
func ValidatePriority(priority int) error {
	if priority < 0 || priority > 65535 {
		return fmt.Errorf("priority must be between 0..65535 (16 bit)")
	}
	return nil
}

// IsDefaultARecord is a predicate to place fetched A domain records into the appropriate bucket
func IsDefaultARecord(record *DomainRecord) bool {
	return record.Name == Ptr && record.Type == AType && record.TTL == DefaultTTL
}

// IsDefaultNSRecord is a predicate to place fetched NS domain records into the appropriate bucket
func IsDefaultNSRecord(record *DomainRecord) bool {
	return record.Name == Ptr && record.Type == NSType && record.TTL == DefaultTTL
}

// IsDisallowed prevents empty NS|SOA record lists from being propagated, which is disallowed
func IsDisallowed(t string, records []*DomainRecord) bool {
	return len(records) == 0 && strings.EqualFold(t, NSType) || strings.EqualFold(t, SOAType)
}

func isSupportedType(recType string) bool {
	for _, t := range supportedTypes {
		if t == recType {
			return true
		}
	}
	return false
}
