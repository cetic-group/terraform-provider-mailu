# mailu_dkim

Reads DKIM and DMARC DNS values exposed by a Mailu domain.

Status: implemented as a read-only data source backed by domain reads.

## Example Usage

```terraform
data "mailu_dkim" "example" {
  domain = "example.com"
}
```

## Schema

### Required

- `domain` (String) Mailu domain name.

### Read-Only

- `id` (String) Same as normalized `domain`.
- `dns_dkim` (String) DKIM DNS value.
- `dns_dmarc` (String) DMARC DNS value.
- `dns_dmarc_report` (String) DMARC report DNS value.
