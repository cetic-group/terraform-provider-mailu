# mailu_token

Manages a Mailu authentication token.

Status: implemented in the extended provider surface.

## Example Usage

```terraform
resource "mailu_token" "admin" {
  email          = "admin@example.com"
  comment        = "Terraform managed token"
  authorized_ips = ["203.0.113.0/24"]
}
```

## Schema

### Required

- `email` (String) User email address that owns the token. Forces replacement.

### Optional

- `comment` (String) Token comment.
- `authorized_ips` (Set of String) Allowed IP addresses or CIDR ranges.

### Read-Only

- `id` (String) Mailu token record identifier.
- `token` (String, Sensitive) Generated token value returned only at creation time.
- `created` (String) API creation timestamp.
- `last_edit` (String) API last-edit timestamp.

## Import

Import using the Mailu token record ID:

```shell
terraform import mailu_token.admin 42
```

Imported tokens do not expose the generated token value because Mailu only returns it during creation.
