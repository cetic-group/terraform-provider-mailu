# Mailu Production Inventory Template

Keep the completed inventory in a private operational location. Do not commit production user lists unless CETIC Group approves it.

## Change Freeze

- Freeze owner:
- Freeze start:
- Freeze end:
- Approved manual exceptions:

## Provider Release

- Provider version:
- Release checksum verified:
- Installed from:
- Operator:

## Domains

| Terraform address | Import ID | Managed | Notes |
| --- | --- | --- | --- |
| `mailu_domain.cetic` | `cetic-group.com` | yes |  |

## Users

| Terraform address | Import ID | Managed | Notes |
| --- | --- | --- | --- |
| `mailu_user.admin` | `admin@cetic-group.com` | yes | Do not set `raw_password` unless rotating. |

## Aliases

| Terraform address | Import ID | Managed | Notes |
| --- | --- | --- | --- |
| `mailu_alias.postmaster` | `postmaster@cetic-group.com` | yes |  |

## Alternative Domains

| Terraform address | Import ID | Managed | Notes |
| --- | --- | --- | --- |
| `mailu_alternative_domain.legacy` | `legacy-cetic-group.com` | no | Example only. |

## Domain Managers

| Terraform address | Import ID | Managed | Notes |
| --- | --- | --- | --- |
| `mailu_domain_manager.admin` | `cetic-group.com/admin@cetic-group.com` | no | Example only. |

## Relays

| Terraform address | Import ID | Managed | Notes |
| --- | --- | --- | --- |
| `mailu_relay.cetic` | `cetic-group.com` | no | SMTP value must not contain credentials. |

## Tokens

| Terraform address | Import ID | Managed | Notes |
| --- | --- | --- | --- |
| `mailu_token.admin` | `42` | no | Token secret value cannot be recovered through Terraform. |

## Unmanaged Objects

Document every Mailu object intentionally left outside Terraform.

| Object type | Object ID | Reason | Owner |
| --- | --- | --- | --- |
|  |  |  |  |
