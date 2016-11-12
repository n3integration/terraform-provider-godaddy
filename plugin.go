package main

import (
	"github.com/hashicorp/terraform/plugin"
)

const (
	base_url = "https://api.godaddy.com/v1"
)

func main() {
	opts := plugin.ServeOpts{
		ProviderFunc: Provider,
	}
	plugin.Serve(&opts)
}
