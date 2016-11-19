package main

import (
	"fmt"
	"log"
)

// Config provides the provider's configuration
type Config struct {
	Key     string
	Secret  string
	BaseURL string
}

// Client returns a new client for accessing GoDaddy.
func (c *Config) Client() (*GoDaddyClient, error) {
	client, err := NewClient(c.BaseURL, c.Key, c.Secret)

	if err != nil {
		return nil, fmt.Errorf("Error setting up client: %s", err)
	}

	log.Print("[INFO] GoDaddy Client configured")

	return client, nil
}
