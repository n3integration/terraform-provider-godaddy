package godaddy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GODADDY_API_KEY", nil),
				Description: "GoDaddy API Key.",
			},

			"secret": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("GODADDY_API_SECRET", nil),
				Description: "GoDaddy API Secret.",
			},

			"baseurl": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https://api.godaddy.com",
				Description: "GoDaddy Base Url(defaults to production).",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"godaddy_domain_record": resourceDomainRecord(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Key:     d.Get("key").(string),
		Secret:  d.Get("secret").(string),
		BaseURL: d.Get("baseurl").(string),
	}

	return config.Client()
}
