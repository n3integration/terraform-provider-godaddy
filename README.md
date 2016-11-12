# terraform-godaddy
[Terraform](https://www.terraform.io/) plugin for managing domain records.

## API Key
In order to leverage the GoDaddy APIs, an [API key](https://developer.godaddy.com/keys/) is required.

## Installation

```bash
glide install && go build
mkdir -p ~/.terraform/plugins && cp terraform-godaddy ~/.terraform/plugins
[ -f ~/.terraformrc ] || cat > ~/.terraformrc <<EOF
providers {
  godaddy = "$HOME/.terraform/plugins/terraform-godaddy"
}
EOF
```

## Usage

```terraform
provider "godaddy" {
  baseurl = "https://api.godaddy.com"
}

resource "godaddy_domain_record" "default" {
  domain = "fancy-domain.com"

  record {
    name = "@"
    type = "A"
    data = "192.168.1.2"
    ttl = 3600
  }

  record {
    name = "@"
    type = "A"
    data = "192.168.1.3"
    ttl = 3600
  }

  record {
    name = "@"
    type = "NS"
    data = "ns7.domains.com"
    ttl = 3600
  }

  record {
    name = "@"
    type = "NS"
    data = "ns6.domains.com"
    ttl = 3600
  }
}
```
