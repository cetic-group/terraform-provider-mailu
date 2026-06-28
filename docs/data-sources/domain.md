# mailu_domain

Reads a Mailu domain.

Status: confirmed by Swagger; not implemented yet.

## Example Usage

```terraform
data "mailu_domain" "example" {
  name = "example.com"
}
```

## Schema

### Required

- `name` (String) Domain name.

### Read-Only

- `id` (String) Domain identifier.
- `comment` (String) Domain comment.
- `managers` (List of String) Domain managers.
- `max_users` (Number) Maximum number of users.
- `max_aliases` (Number) Maximum number of aliases.
- `max_quota_bytes` (Number) Maximum mailbox quota in bytes.
- `signup_enabled` (Boolean) Whether signup is enabled.
- `alternatives` (List of String) Alternative domain names.
- `dns_autoconfig` (List of String) Autoconfiguration DNS records.
- `dns_mx` (String) MX DNS value.
- `dns_spf` (String) SPF DNS value.
- `dns_dkim` (String) DKIM DNS value.
- `dns_dmarc` (String) DMARC DNS value.
- `dns_dmarc_report` (String) DMARC report DNS value.
- `dns_tlsa` (List of String) TLSA DNS values.
