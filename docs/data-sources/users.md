---
page_title: "mailu_users Data Source - terraform-provider-mailu"
subcategory: ""
description: |-
  Lists all Mailu mailbox users. Useful for inventorying existing objects when generating Terraform import blocks.
---

# mailu_users (Data Source)

Lists all Mailu mailbox users. Useful for inventorying existing objects when generating Terraform import blocks.

## Example Usage

```terraform
data "mailu_users" "all" {}

output "user_emails" {
  value = [for u in data.mailu_users.all.users : u.email]
}
```

## Schema

### Read-Only

- `users` (Attributes List) (see [below for nested schema](#nestedatt--users))

<a id="nestedatt--users"></a>
### Nested Schema for `users`

Read-Only:

- `email` (String)
- `comment` (String)
- `quota_bytes` (Number)
- `enabled` (Boolean)
- `global_admin` (Boolean)
- `displayed_name` (String)
