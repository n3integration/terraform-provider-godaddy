# terraform-godaddy

## Installation

```bash
mkdir -p ~/.terraform/plugins
[ -f ~/.terraformrc ] || cat > ~/.terraformrc <<EOF
providers {
  godaddy = "$HOME/.terraform/plugins/terraform-godaddy"
}
EOF
```

## API Key
In order to leverage the GoDaddy APIs, an [API key](https://developer.godaddy.com/keys/) is required.
