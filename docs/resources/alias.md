# mailu_alias

Manages a Mailu alias.

Status: confirmed by Swagger; not implemented yet.

## Example Usage

```terraform
resource "mailu_alias" "postmaster" {
  email       = "postmaster@example.com"
  destination = ["admin@example.com"]
}
```

## Schema

### Required

- `email` (String) Alias email address.

### Optional

- `destination` (Set of String) Destination email addresses.
- `comment` (String) Alias comment.
- `wildcard` (Boolean) Enable SQL LIKE wildcard syntax.

### Read-Only

- `id` (String) Alias identifier.

## Import

Import using the alias email address:

```shell
terraform import mailu_alias.postmaster postmaster@example.com
```
