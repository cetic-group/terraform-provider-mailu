# mailu_domain

Manages a Mailu domain.

Status: confirmed by Swagger; not implemented yet.

## Example Usage

```terraform
resource "mailu_domain" "example" {
  name = "example.com"
}
```

## Schema

### Required

- `name` (String) Domain name.

### Optional

- `comment` (String) Domain comment.
- `max_users` (Number) Maximum number of users.
- `max_aliases` (Number) Maximum number of aliases.
- `max_quota_bytes` (Number) Maximum mailbox quota in bytes.
- `signup_enabled` (Boolean) Whether signup is enabled.
- `alternatives` (List of String) Alternative domain names.

### Read-Only

- `id` (String) Domain identifier.
- `managers` (List of String) Domain managers.
- `dns_autoconfig` (List of String) Autoconfiguration DNS records.
- `dns_mx` (String) MX DNS value.
- `dns_spf` (String) SPF DNS value.
- `dns_dkim` (String) DKIM DNS value.
- `dns_dmarc` (String) DMARC DNS value.
- `dns_dmarc_report` (String) DMARC report DNS value.
- `dns_tlsa` (List of String) TLSA DNS values.

## Import

Import using the domain name:

```shell
terraform import mailu_domain.example example.com
```
