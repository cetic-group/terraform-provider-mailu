# mailu_alternative_domain

Manages a Mailu alternative domain.

Status: implemented in the extended provider surface.

## Example Usage

```terraform
resource "mailu_alternative_domain" "example" {
  name   = "example.net"
  domain = "example.com"
}
```

## Schema

### Required

- `name` (String) Alternative domain name. Forces replacement.
- `domain` (String) Parent Mailu domain name. Forces replacement.

### Read-Only

- `id` (String) Alternative domain identifier. Same as normalized `name`.

## Import

Import using the alternative domain name:

```shell
terraform import mailu_alternative_domain.example example.net
```
