package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

const (
	// BaseURL points to the GoDaddy API's base
	BaseURL = "https://api.godaddy.com/v1"
)

func main() {
	opts := &plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return Provider()
		},
	}
	plugin.Serve(opts)
}
