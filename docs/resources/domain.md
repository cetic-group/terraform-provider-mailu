# mailu_domain

Manages a Mailu domain.

Status: implemented in the MVP provider surface.

## Example Usage

```terraform
resource "mailu_domain" "example" {
  name            = "example.com"
  comment         = "Managed by Terraform"
  max_users       = 25
  max_aliases     = 50
  max_quota_bytes = 10737418240
  signup_enabled  = false
}
```

## Schema

### Required

- `name` (String) Domain name. Forces replacement.

### Optional

- `comment` (String) Domain comment.
- `max_users` (Number) Maximum number of users.
- `max_aliases` (Number) Maximum number of aliases.
- `max_quota_bytes` (Number) Maximum mailbox quota in bytes.
- `signup_enabled` (Boolean) Whether signup is enabled.
- `alternatives` (Set of String) Alternative domain names.

### Read-Only

- `id` (String) Domain identifier. Same as normalized `name`.
- `managers` (Set of String) Domain managers.
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
