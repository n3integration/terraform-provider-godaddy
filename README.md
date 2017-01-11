# terraform-godaddy
[Terraform](https://www.terraform.io/) plugin for managing domain records

[ ![Codeship Status for n3integration/terraform-godaddy](https://app.codeship.com/projects/29e8c490-8b5d-0134-914d-3e63d62140d1/status?branch=master)](https://app.codeship.com/projects/184616)

<dl>
  <dt>Terraform v0.7.x</dt>
  <dd>https://github.com/n3integration/terraform-godaddy/releases/tag/v1.0.0</dd>
  <dt>Terraform v0.8.x</dt>
  <dd>https://github.com/n3integration/terraform-godaddy/releases/tag/v1.1.1</dd>
<dl>

## Installation

```bash
bash <(curl -s https://raw.githubusercontent.com/n3integration/terraform-godaddy/master/install.sh)
```

## API Key
In order to leverage the GoDaddy APIs, an [API key](https://developer.godaddy.com/keys/) is required. The key pair can be optionally stored in environment variables.

```bash
export GD_KEY=abc
export GD_SECRET=123
```

## Provider

If `key` and `secret` aren't provided under the `godaddy` `provider`, they are expected to be exposed as environment variables: `GD_KEY` and `GD_SECRET`.

```terraform
provider "godaddy" {
  key = "abc"
  secret = "123"
}
```

## Domain Record Resource
A `godaddy_domain_record` resource requires a `domain`. If the domain is not registered under the account that owns the key, an optional `customer` number can be specified. 
Additionally, one or more `record` instances are required. For each `record`, the `name`, `type`, and `data` attributes are required. The available types include:

* A
* AAAA
* CNAME
* NS
* SOA
* TXT

```terraform
resource "godaddy_domain_record" "default" {
  domain = "fancy-domain.com"
  customer = "1234"

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
    name = "www"
    type = "CNAME"
    data = "fancy.github.io"
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

## License

Copyright 2016 n3integration@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
