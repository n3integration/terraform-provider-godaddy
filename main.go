package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"



	"github.com/kolikons/terraform-provider-godaddy/godaddy"
)

const (
	// BaseURL points to the GoDaddy API's base
	BaseURL = "https://api.godaddy.com/v1"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return godaddy.Provider()
		},
	})
}
