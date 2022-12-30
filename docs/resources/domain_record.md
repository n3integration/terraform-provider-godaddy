---
page_title: "godaddy_domain_record Resource - terraform-provider-godaddy"
subcategory: "infrastructure"
<<<<<<< Updated upstream

=======
description: |-
  
>>>>>>> Stashed changes
---

# godaddy_domain_record (Resource)

## Schema

### Required

- **domain** (String)

### Optional

<<<<<<< Updated upstream
- **addresses** (List of String)
- **customer** (String)
- **id** (String) The ID of this resource.
- **nameservers** (List of String)
- **record** (Block Set) (see [below for nested schema](#nestedblock--record))
=======
- `addresses` (List of String) IP Addresses.
- `customer` (String) Customer ID (required if you are a reseller managing a domain purchased outside the scope of your reseller account).
- `nameservers` (List of String)
- `record` (Block Set) (see [below for nested schema](#nestedblock--record))

### Read-Only

- `id` (String) The ID of this resource.
>>>>>>> Stashed changes

<a id="nestedblock--record"></a>
### Nested Schema for `record`

Required:

- **data** (String)
- **name** (String)
- **type** (String)

Optional:

- **port** (Number)
- **priority** (Number)
- **protocol** (String)
- **service** (String)
- **ttl** (Number)
- **weight** (Number)
