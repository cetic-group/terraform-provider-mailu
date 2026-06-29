# DNS Integration

This provider manages Mailu state. It does not own DNS zones.

Use Mailu resources and data sources to expose mail DNS values, then pass those values to the DNS provider that owns the zone.

## Ownership Boundary

| Area | Owner |
| --- | --- |
| Mailu domains, users, aliases, relays, managers, tokens | `cetic-group/mailu` provider |
| DNS zones and DNS records | DNS provider, for example IONOS or another DNS provider |
| MTA-STS HTTPS policy file | Web/static hosting or reverse proxy configuration |

Do not create DNS-specific resources in this provider. DNS providers have different schemas, validation rules, TTL behavior, and zone ownership models.

## Record Classes

| Record | Priority | Source | Notes |
| --- | --- | --- | --- |
| MX | Mandatory | `mailu_domain.dns_mx` or `data.mailu_domain.dns_mx` | Required for inbound mail delivery. |
| SPF TXT | Mandatory | `mailu_domain.dns_spf` or `data.mailu_domain.dns_spf` | Required for outbound deliverability. |
| DKIM TXT | Mandatory after DKIM exists | `data.mailu_dkim.dns_dkim` | Required for modern deliverability. Generate/rotate keys explicitly in Mailu, then read DNS through Terraform. |
| DMARC TXT | Recommended | `data.mailu_dkim.dns_dmarc` | Start with a monitoring policy before enforcing reject/quarantine. |
| DMARC report TXT | Optional | `data.mailu_dkim.dns_dmarc_report` | Present when Mailu exposes a report value. |
| Autoconfig CNAME/SRV | Recommended | `mailu_domain.dns_autoconfig` | Improves mail client setup. |
| TLSA | Optional | `mailu_domain.dns_tlsa` | Use when DNSSEC/DANE is operational. |
| MTA-STS TXT | Recommended | Manual Terraform input | Requires `_mta-sts.<domain>` TXT and a served policy file. |
| TLS reporting TXT | Recommended | Manual Terraform input | `_smtp._tls.<domain>` TXT for delivery reports. |

## Pattern

1. Create or import the Mailu domain with `mailu_domain`.
2. Read DNS values from `mailu_domain` and `mailu_dkim`.
3. Convert Mailu DNS zone-line strings to the schema required by the DNS provider.
4. Apply DNS records with the DNS provider.
5. Keep DNS credentials outside this provider and outside source files.

## DKIM

`mailu_dkim` is read-only. It reads DKIM and DMARC values from the Mailu domain response.

The provider intentionally does not model DKIM generation as a Terraform resource. `POST /domain/{domain}/dkim` rotates keys and is an operator action, not normal CRUD state. Run key generation explicitly, then apply DNS records from the new `data.mailu_dkim` values.

## MTA-STS

Terraform can manage DNS records for MTA-STS, but Mailu does not expose the HTTPS policy file through this provider.

Typical DNS records:

```text
_mta-sts.example.com.  TXT  "v=STSv1; id=2026062901"
_smtp._tls.example.com. TXT  "v=TLSRPTv1; rua=mailto:postmaster@example.com"
```

The matching policy must be served at:

```text
https://mta-sts.example.com/.well-known/mta-sts.txt
```

Example policy:

```text
version: STSv1
mode: enforce
mx: mail.example.com
max_age: 604800
```

Use `mode: testing` before enforcing a production domain.

## IONOS

For CETIC Group IONOS DNS automation, keep the IONOS API key outside Git. The existing Mailu deployment already documents the expected env file shape in `mailu/mailu-data/ionos.env.example`.

The IONOS-oriented example under `examples/dns/ionos` produces a normalized list of records and an API-style payload that can be mapped to the selected IONOS automation layer.
