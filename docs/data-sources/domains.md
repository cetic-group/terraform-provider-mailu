---
page_title: "mailu_domains Data Source - terraform-provider-mailu"
subcategory: ""
description: |-
  Lists all Mailu domains. Useful for inventorying existing objects when generating Terraform import blocks.
---

# mailu_domains (Data Source)

Lists all Mailu domains. Useful for inventorying existing objects when generating Terraform import blocks.

## Example Usage

```terraform
data "mailu_domains" "all" {}

output "domain_names" {
  value = [for d in data.mailu_domains.all.domains : d.name]
}
```

## Schema

### Read-Only

- `domains` (Attributes List) (see [below for nested schema](#nestedatt--domains))

<a id="nestedatt--domains"></a>
### Nested Schema for `domains`

Read-Only:

- `name` (String)
- `comment` (String)
- `max_users` (Number)
- `max_aliases` (Number)
- `max_quota_bytes` (Number)
- `signup_enabled` (Boolean)
- `alternatives` (Set of String)
