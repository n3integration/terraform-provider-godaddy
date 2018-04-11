package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	headerAccept        = "Accept"
	headerAuthorization = "Authorization"
	headerContent       = "Content-Type"
	headerCustomerID    = "X-Shopper-Id"
	mediaTypeJSON       = "application/json"
	pathDomainRecords   = "%s/v1/domains/%s/records"
	pathDomains         = "%s/v1/domains/%s"
	rateLimit           = 1 * time.Second
)

// GoDaddyClient is a GoDaddy API client
type GoDaddyClient struct {
	baseURL string
	key     string
	secret  string
	client  *http.Client
}

// rateLimitedTransport throttles API calls to GoDaddy. It appears that
// the rate limit is 60 requests per minute, which can be throttled and
// enforced at a maximum of one request/second.
type rateLimitedTransport struct {
	delegate http.RoundTripper
	throttle time.Time
	sync.Mutex
}

func (t *rateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	t.Lock()
	defer t.Unlock()

	if t.throttle.After(time.Now()) {
		delta := t.throttle.Sub(time.Now())
		time.Sleep(delta)
	}

	t.throttle = time.Now().Add(rateLimit)
	return t.delegate.RoundTrip(req)
}

// NewClient constructs a new GoDaddy API client or an error if the supplied
// input is invalid.
func NewClient(baseURL, key, secret string) (*GoDaddyClient, error) {
	baseURL, err := formatURL(baseURL)
	if err != nil {
		return nil, err
	}

	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10 * time.Second,
	}

	return &GoDaddyClient{
		baseURL: baseURL,
		key:     strings.TrimSpace(key),
		secret:  strings.TrimSpace(secret),
		client: &http.Client{
			Timeout: time.Second * 30,
			Transport: &rateLimitedTransport{
				delegate: netTransport,
				throttle: time.Now().Add(-(rateLimit)),
			},
		},
	}, nil
}

// GetDomains fetches the details for the provided domain
func (c *GoDaddyClient) GetDomains(customerID string) ([]Domain, error) {
	url := fmt.Sprintf(pathDomains, c.baseURL, "")

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	d := []Domain{}
	if err := c.execute(customerID, req, &d); err != nil {
		return nil, err
	}

	return d, nil
}

// GetDomain fetches the details for the provided domain
func (c *GoDaddyClient) GetDomain(customerID, domain string) (*Domain, error) {
	url := fmt.Sprintf(pathDomains, c.baseURL, domain)

	req, err := http.NewRequest(http.MethodGet, url, nil)
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
func (c *GoDaddyClient) GetDomainRecords(customerID, domain string) ([]*DomainRecord, error) {
	url := fmt.Sprintf(pathDomainRecords, c.baseURL, domain)

	req, err := http.NewRequest(http.MethodGet, url, nil)
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
func (c *GoDaddyClient) UpdateDomainRecords(customerID, domain string, records []*DomainRecord) error {
	msg, err := json.Marshal(records)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(pathDomainRecords, c.baseURL, domain)
	method := http.MethodPut

	req, err := http.NewRequest(method, url, bytes.NewBuffer(msg))
	if err != nil {
		return err
	}

	return c.execute(customerID, req, nil)
}

func (c *GoDaddyClient) execute(customerID string, req *http.Request, result interface{}) error {
	if len(strings.TrimSpace(customerID)) > 0 {
		req.Header.Set(headerCustomerID, customerID)
	}

	req.Header.Set(headerAccept, mediaTypeJSON)
	req.Header.Set(headerContent, mediaTypeJSON)
	req.Header.Set(headerAuthorization, fmt.Sprintf("sso-key %s:%s", c.key, c.secret))

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err = validate(resp); err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if result == nil {
		return nil
	} else if err = json.Unmarshal(body, result); err != nil {
		return err
	}

	return nil
}

func validate(resp *http.Response) error {
	if resp.StatusCode < http.StatusBadRequest {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var errResp = struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{}

	if err := json.Unmarshal(body, &errResp); err != nil {
		return err
	}

	return fmt.Errorf("[%d:%s] %s %s", resp.StatusCode, errResp.Code, errResp.Message, body)
}

func formatURL(baseURL string) (string, error) {
	url, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	if url.Host == "" || url.Scheme == "" {
		return "", fmt.Errorf("invalid baseUrl. expected format: scheme://host")
	}

	return fmt.Sprintf("%s://%s", url.Scheme, url.Host), err
}
