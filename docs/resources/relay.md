# mailu_relay

Manages a Mailu relay.

Status: implemented in the extended provider surface.

## Example Usage

```terraform
resource "mailu_relay" "example" {
  name    = "example.com"
  smtp    = "smtp.example.net"
  comment = "Managed by Terraform"
}
```

## Schema

### Required

- `name` (String) Relayed domain name. Forces replacement.

### Optional

- `smtp` (String) Remote SMTP host.
- `comment` (String) Relay comment.

### Read-Only

- `id` (String) Relay identifier. Same as normalized `name`.

## Import

Import using the relay name:

```shell
terraform import mailu_relay.example example.com
```
