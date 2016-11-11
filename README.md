# terraform-godaddy

## Installation

```bash
cd ~ && mkdir -p .terraform/plugins
[ -f ~/.terraformrc ] || cat <-EOF>> ~/.terraformrc
providers {
  godaddy = ~/.terraform/plugins/godaddy
}
EOF
```

## API Key
In order to leverage the GoDaddy APIs, an [API key](https://developer.godaddy.com/keys/) is required.
