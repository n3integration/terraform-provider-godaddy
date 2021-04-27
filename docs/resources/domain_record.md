---
page_title: "godaddy_domain_record Resource - terraform-provider-godaddy"
subcategory: "infrastructure"

---

# godaddy_domain_record (Resource)

## Schema

### Required

- **domain** (String)

### Optional

- **addresses** (List of String)
- **customer** (String)
- **id** (String) The ID of this resource.
- **nameservers** (List of String)
- **record** (Block Set) (see [below for nested schema](#nestedblock--record))

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
