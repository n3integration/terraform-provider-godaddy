---
page_title: "GoDaddy Provider"
subcategory: "infrastructure"
description: |-
  
---

# GoDaddy Provider

The GoDaddy Provider supports Terraform's `import` feature to make it easier to deal with existing records that are supplied by GoDaddy by default.

#### Import Example

To import any pre-existing GoDaddy resource data into your local Terraform state file, Terraform should be invoked using the `import` command with the fully-qualified resource name and domain name as command arguments. For example:

```bash
terraform import godaddy_domain_record.mydomain mydomain.com
```

## Schema

### Optional

- **baseurl** (String) GoDaddy Base URL(defaults to production).
- **key** (String) GoDaddy API Key.
- **secret** (String) GoDaddy API Secret.
