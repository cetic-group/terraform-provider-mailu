# mailu_dkim

Reads DKIM metadata for a Mailu domain.

Status: partial. Swagger exposes DKIM generation and domain DNS output, but no dedicated DKIM read endpoint.

## Example Usage

```terraform
data "mailu_dkim" "example" {
  domain = "example.com"
}
```

## Schema

### Required

- `domain` (String) Mailu domain.

### Read-Only

- `selector` (String) DKIM selector.
- `public_key` (String) DKIM public key.
- `dns_name` (String) DNS record name.
- `dns_value` (String) DNS record value.
