package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	rateLimit           = 1 * time.Second
)

// Client is a GoDaddy API client
type Client struct {
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
func NewClient(baseURL, key, secret string) (*Client, error) {
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

	return &Client{
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

func (c *Client) execute(customerID string, req *http.Request, result interface{}) error {
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

	buffer := new(bytes.Buffer)
	body, err := ioutil.ReadAll(io.TeeReader(resp.Body, buffer))
	log.Printf("%s %s", resp.Status, buffer)
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

	buffer := new(bytes.Buffer)
	body, err := ioutil.ReadAll(io.TeeReader(resp.Body, buffer))
	log.Println(buffer)
	if err != nil {
		return err
	}

	var errResp = struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Fields  []struct {
			Code        string `json:"code"`
			Message     string `json:"message"`
			Path        string `json:"path"`
			PathRelated string `json:"pathRelated"`
		} `json:"fields"`
	}{}

	if err := json.Unmarshal(body, &errResp); err != nil {
		return err
	}

	if len(errResp.Fields) == 0 {
		return fmt.Errorf("[%d:%s] %s", resp.StatusCode, errResp.Code, errResp.Message)
	}

	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("[%d:%s] %s (", resp.StatusCode, errResp.Code, errResp.Message))
	for i, field := range errResp.Fields {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(fmt.Sprintf("%s [%s]: %s", field.Path, field.Code, field.Message))
	}
	b.WriteString(")")
	return fmt.Errorf("%s", b.String())
}

func formatURL(base string) (string, error) {
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	if baseURL.Host == "" || baseURL.Scheme == "" {
		return "", fmt.Errorf("invalid baseUrl. expected format: scheme://host")
	}

	return fmt.Sprintf("%s://%s", baseURL.Scheme, baseURL.Host), err
}
