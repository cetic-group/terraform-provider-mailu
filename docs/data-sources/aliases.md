---
page_title: "mailu_aliases Data Source - terraform-provider-mailu"
subcategory: ""
description: |-
  Lists all Mailu aliases. Useful for inventorying existing objects when generating Terraform import blocks.
---

# mailu_aliases (Data Source)

Lists all Mailu aliases. Useful for inventorying existing objects when generating Terraform import blocks.

## Example Usage

```terraform
data "mailu_aliases" "all" {}

output "alias_emails" {
  value = [for a in data.mailu_aliases.all.aliases : a.email]
}
```

## Schema

### Read-Only

- `aliases` (Attributes List) (see [below for nested schema](#nestedatt--aliases))

<a id="nestedatt--aliases"></a>
### Nested Schema for `aliases`

Read-Only:

- `email` (String)
- `destination` (Set of String)
- `comment` (String)
- `wildcard` (Boolean)
