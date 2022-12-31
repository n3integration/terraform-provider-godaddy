---
page_title: "godaddy_domain_record Resource - terraform-provider-godaddy"
subcategory: "infrastructure"
description: |-
  
---

# godaddy_domain_record (Resource)

## Schema

### Required

- `domain` (String)

### Optional

- `addresses` (List of String) IP Addresses.
- `customer` (String) Customer ID (required if you are a reseller managing a domain purchased outside the scope of your reseller account).
- `nameservers` (List of String)
- `record` (Block Set) (see [below for nested schema](#nestedblock--record))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--record"></a>
### Nested Schema for `record`

Required:

- `data` (String)
- `name` (String)
- `type` (String)

Optional:

- `port` (Number)
- `priority` (Number)
- `protocol` (String)
- `service` (String)
- `ttl` (Number)
- `weight` (Number)
