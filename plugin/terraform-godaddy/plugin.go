package main

import (
	"github.com/hashicorp/terraform/plugin"
)

const (
	// BaseURL points to the GoDaddy API's base
	BaseURL = "https://api.godaddy.com/v1"
)

func main() {
	opts := plugin.ServeOpts{
		ProviderFunc: Provider,
	}
	plugin.Serve(&opts)
}
