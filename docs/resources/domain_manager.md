# mailu_domain_manager

Manages a Mailu domain manager assignment.

Status: implemented in the extended provider surface.

## Example Usage

```terraform
resource "mailu_domain_manager" "admin" {
  domain     = "example.com"
  user_email = "admin@example.com"
}
```

## Schema

### Required

- `domain` (String) Domain managed by the user. Forces replacement.
- `user_email` (String) Manager user email address. Forces replacement.

### Read-Only

- `id` (String) Assignment identifier in `<domain>/<email>` format.

## Import

Import using `<domain>/<email>`:

```shell
terraform import mailu_domain_manager.admin example.com/admin@example.com
```
