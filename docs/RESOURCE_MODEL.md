# Resource Model

This document tracks the Terraform model for the `cetic-group/mailu` provider.

Status: MVP and phase 6 extended resource models implemented.

Resource modeling decisions require architecture, QA, and security review. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Provider

```hcl
provider "mailu" {
  endpoint = "https://mail.cetic-group.com/api/v1"
  token    = var.mailu_api_token
}
```

| Attribute | Kind | Notes |
| --- | --- | --- |
| `endpoint` | optional | Mailu API base URL. Can be set with `MAILU_ENDPOINT`. |
| `token` | optional, sensitive | Mailu API token. Can be set with `MAILU_API_TOKEN`. |

## Common Rules

Normalization:

- Domain names are normalized to lowercase and trimmed.
- Email addresses are normalized to lowercase and trimmed.
- List/set attributes containing email addresses or domains are normalized item-by-item.

Terraform IDs:

- Use stable natural identifiers.
- Domain ID: domain name.
- User ID: full email address.
- Alias ID: full alias email address.
- Alternative domain ID: alternative domain name.
- Domain manager ID: `<domain>/<email>`.
- Relay ID: relay name.
- Token ID: Mailu token record ID.

Delete behavior:

- MVP resources map Terraform delete to Mailu `DELETE`.
- Runtime validation confirmed hard delete behavior for domains, users, and aliases.
- Reads after delete return `404`.

Drift behavior:

- Managed mutable fields are compared against Mailu read responses.
- Computed DNS and usage fields are refreshed from Mailu and do not force updates.
- API-returned password hashes are ignored and never compared with `raw_password`.

Sensitive handling:

- `token` and `raw_password` are sensitive.
- `UserGet.password` is a hash returned by Mailu and must not be exposed as a Terraform attribute.
- Diagnostics and logs must redact tokens, raw passwords, generated token values, and password hashes.

## `mailu_domain`

Status: MVP schema frozen.

Represents a Mailu domain.

Endpoint mapping:

- Create/list: `POST /domain`, `GET /domain`
- Read/update/delete: `GET /domain/{domain}`, `PATCH /domain/{domain}`, `DELETE /domain/{domain}`

Import ID:

```text
example.com
```

Schema:

| Attribute | Kind | Replacement | Notes |
| --- | --- | --- | --- |
| `id` | computed | no | Same as normalized `name`. |
| `name` | required | yes | Domain name. Mailu does not expose rename. |
| `comment` | optional | no | Domain comment. |
| `max_users` | optional | no | Maximum users, Mailu default may be `-1`. |
| `max_aliases` | optional | no | Maximum aliases, Mailu default may be `-1`. |
| `max_quota_bytes` | optional | no | Maximum mailbox quota in bytes. |
| `signup_enabled` | optional | no | Whether signup is enabled. |
| `alternatives` | optional | no | Set of alternative domain names. |
| `managers` | computed | no | Domain managers. |
| `dns_autoconfig` | computed | no | Autoconfiguration DNS records. |
| `dns_mx` | computed | no | MX DNS value. |
| `dns_spf` | computed | no | SPF DNS value. |
| `dns_dkim` | computed | no | DKIM DNS value. |
| `dns_dmarc` | computed | no | DMARC DNS value. |
| `dns_dmarc_report` | computed | no | DMARC report DNS value. |
| `dns_tlsa` | computed | no | TLSA DNS values. |

## `mailu_user`

Status: MVP schema frozen.

Represents a mailbox user.

Endpoint mapping:

- Create/list: `POST /user`, `GET /user`
- Read/update/delete: `GET /user/{email}`, `PATCH /user/{email}`, `DELETE /user/{email}`
- Domain listing: `GET /domain/{domain}/users`

Import ID:

```text
user@example.com
```

Schema:

| Attribute | Kind | Replacement | Notes |
| --- | --- | --- | --- |
| `id` | computed | no | Same as normalized `email`. |
| `email` | required | yes | Full email address. Mailu does not expose rename. |
| `raw_password` | required on create, optional on update, sensitive | no | Write-only input mapped to Mailu `raw_password`. |
| `comment` | optional | no | User comment. |
| `quota_bytes` | optional | no | Mailbox quota in bytes. |
| `quota_bytes_used` | computed | no | Used quota in bytes. |
| `global_admin` | optional | no | Whether the user is a global Mailu admin. |
| `enabled` | optional | no | Whether the user is enabled. |
| `change_pw_next_login` | optional | no | Force password change at next login. |
| `enable_imap` | optional | no | Allow IMAP access. |
| `enable_pop` | optional | no | Allow POP3 access. |
| `allow_spoofing` | optional | no | Allow sender spoofing. |
| `forward_enabled` | optional | no | Enable forwarding. |
| `forward_destination` | optional | no | Set of forward destination email addresses. |
| `forward_keep` | optional | no | Keep a copy when forwarding. |
| `reply_enabled` | optional | no | Enable automatic replies. |
| `reply_subject` | optional | no | Automatic reply subject. |
| `reply_body` | optional, sensitive | no | Automatic reply body can contain personal data. |
| `reply_startdate` | optional | no | Date string in `YYYY-MM-DD` format. |
| `reply_enddate` | optional | no | Date string in `YYYY-MM-DD` format. |
| `displayed_name` | optional | no | Display name. |
| `spam_enabled` | optional | no | Enable spam filtering. |
| `spam_mark_as_read` | optional | no | Mark spam as read. |
| `spam_threshold` | optional | no | Spam threshold. |

Password update behavior:

- Setting or changing `raw_password` sends `raw_password` in `PATCH /user/{email}`.
- `raw_password` is never read back from Mailu.
- The API-returned `password` hash is ignored and not stored as a normal schema field.

## `mailu_alias`

Status: MVP schema frozen.

Represents a Mailu alias.

Endpoint mapping:

- Create/list: `POST /alias`, `GET /alias`
- Read/update/delete: `GET /alias/{alias}`, `PATCH /alias/{alias}`, `DELETE /alias/{alias}`
- Domain filter: `GET /alias/destination/{domain}`

Import ID:

```text
alias@example.com
```

Schema:

| Attribute | Kind | Replacement | Notes |
| --- | --- | --- | --- |
| `id` | computed | no | Same as normalized `email`. |
| `email` | required | yes | Full alias email address. Mailu does not expose rename. |
| `destination` | optional | no | Set of destination email addresses. |
| `comment` | optional | no | Alias comment. |
| `wildcard` | optional | no | Enable SQL LIKE wildcard syntax. |

## Data Sources

Implemented data sources:

- `mailu_domain`
- `mailu_user`
- `mailu_dkim`

Data source IDs use the same natural IDs as resources.

## `mailu_alternative_domain`

Status: implemented.

Represents an alternative Mailu domain.

Endpoint mapping:

- Create/list: `POST /alternative`, `GET /alternative`
- Read/delete: `GET /alternative/{alt}`, `DELETE /alternative/{alt}`

Import ID:

```text
example.net
```

Schema:

| Attribute | Kind | Replacement | Notes |
| --- | --- | --- | --- |
| `id` | computed | no | Same as normalized `name`. |
| `name` | required | yes | Alternative domain name. |
| `domain` | required | yes | Parent Mailu domain name. |

## `mailu_domain_manager`

Status: implemented.

Represents a manager assignment on a Mailu domain.

Endpoint mapping:

- Create/list: `POST /domain/{domain}/manager`, `GET /domain/{domain}/manager`
- Read/delete: `GET /domain/{domain}/manager/{email}`, `DELETE /domain/{domain}/manager/{email}`

Import ID:

```text
example.com/admin@example.com
```

Schema:

| Attribute | Kind | Replacement | Notes |
| --- | --- | --- | --- |
| `id` | computed | no | `<domain>/<email>`. |
| `domain` | required | yes | Domain managed by the user. |
| `user_email` | required | yes | Manager user email address. |

## `mailu_relay`

Status: implemented.

Represents a Mailu relay.

Endpoint mapping:

- Create/list: `POST /relay`, `GET /relay`
- Read/update/delete: `GET /relay/{name}`, `PATCH /relay/{name}`, `DELETE /relay/{name}`

Import ID:

```text
example.com
```

Schema:

| Attribute | Kind | Replacement | Notes |
| --- | --- | --- | --- |
| `id` | computed | no | Same as normalized `name`. |
| `name` | required | yes | Relayed domain name. |
| `smtp` | optional | no | Remote SMTP host. |
| `comment` | optional | no | Relay comment. |

## `mailu_token`

Status: implemented.

Represents a Mailu authentication token.

Endpoint mapping:

- Create/list: `POST /token`, `GET /token`
- Read/update/delete: `GET /token/{token_id}`, `PATCH /token/{token_id}`, `DELETE /token/{token_id}`

Import ID:

```text
42
```

Schema:

| Attribute | Kind | Replacement | Notes |
| --- | --- | --- | --- |
| `id` | computed | no | Mailu token record ID. |
| `email` | required | yes | User that owns the token. |
| `comment` | optional | no | Token comment. |
| `authorized_ips` | optional | no | Allowed IP addresses or networks. |
| `token` | computed, sensitive | no | Generated token value returned only during creation. |
| `created` | computed | no | API creation timestamp. |
| `last_edit` | computed | no | API last-edit timestamp. |

Security behavior:

- `token` is marked sensitive and is never logged by provider diagnostics.
- Imported tokens cannot recover the generated token value because Mailu does not return it after creation.

## `mailu_dkim`

Status: implemented as a data source.

Reads DKIM and DMARC DNS values from `GET /domain/{domain}`.

No Terraform resource generates DKIM keys because `POST /domain/{domain}/dkim` is an action with destructive rotation semantics. Key generation should remain an explicit operator action until a safe lifecycle model is designed.

## Acceptance Test Fixtures

Acceptance tests must be opt-in with `TF_ACC=1`.

Required environment variables:

- `MAILU_ENDPOINT`
- `MAILU_API_TOKEN`
- `MAILU_ACC_DOMAIN`

Fixture rules:

- Create only temporary domains matching `tf-acc-*.<MAILU_ACC_DOMAIN>`.
- Create users and aliases only under the temporary domain.
- Clean up aliases before users and users before domains.
- Verify cleanup with reads that return `404`.
- Refuse to run when `MAILU_ACC_DOMAIN` is empty.

## Deferred Or Unsupported

- `mailu_fetchmail`: no Swagger endpoint found.
- `mailu_server_info`: no Swagger endpoint found.
- Standalone DKIM generation resource: `POST /domain/{domain}/dkim` rotates keys and is intentionally not modeled as CRUD.
