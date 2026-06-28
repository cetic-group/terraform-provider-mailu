# Resource Model

This document describes the intended Terraform model for the `cetic-group/mailu` provider.

Resource modeling decisions require architecture, QA, and security review. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Provider

```hcl
provider "mailu" {
  endpoint = "https://mail.cetic-group.com/api/v1"
  token    = var.mailu_api_token
}
```

Attributes:

- `endpoint`: Mailu API base URL.
- `token`: sensitive API token.

Environment variables:

- `MAILU_ENDPOINT`
- `MAILU_API_TOKEN`

## `mailu_domain`

Status: confirmed by Swagger; runtime validation pending.

Represents a Mailu domain.

Confirmed attributes:

- `name`: domain name, required, force replacement if Mailu cannot rename domains.
- `comment`: optional.
- `max_users`: optional.
- `max_aliases`: optional.
- `max_quota_bytes`: optional.
- `signup_enabled`: optional.
- `alternatives`: optional list of alternative domain names.

Computed attributes:

- `managers`
- `dns_autoconfig`
- `dns_mx`
- `dns_spf`
- `dns_dkim`
- `dns_dmarc`
- `dns_dmarc_report`
- `dns_tlsa`

Import ID:

```text
example.com
```

## `mailu_user`

Status: confirmed by Swagger; runtime validation pending.

Represents a mailbox user.

Confirmed attributes:

- `email`: required, force replacement if Mailu cannot rename users.
- `raw_password`: required on create, optional on update, sensitive, write-only.
- `comment`: optional.
- `quota_bytes`: optional.
- `enabled`: optional, default `true`.
- `global_admin`: optional, default `false`.
- `change_pw_next_login`: optional.
- `enable_imap`: optional.
- `enable_pop`: optional.
- `allow_spoofing`: optional.
- `forward_enabled`: optional.
- `forward_destination`: optional list of email addresses.
- `forward_keep`: optional.
- `reply_enabled`: optional.
- `reply_subject`: optional.
- `reply_body`: optional.
- `reply_startdate`: optional date string.
- `reply_enddate`: optional date string.
- `displayed_name`: optional.
- `spam_enabled`: optional.
- `spam_mark_as_read`: optional.
- `spam_threshold`: optional.

Computed attributes:

- `quota_bytes_used`

Import ID:

```text
user@example.com
```

## `mailu_alias`

Status: confirmed by Swagger; runtime validation pending.

Represents a Mailu alias.

Confirmed attributes:

- `email`: required.
- `destination`: optional list/set of destination email addresses.
- `comment`: optional.
- `wildcard`: optional.

Import ID:

```text
alias@example.com
```

## Password Handling

`raw_password` must be treated as sensitive write-only input. Mailu returns `UserGet.password`, which is a password hash. The provider must not expose this hash as a normal attribute and must not compare it with `raw_password`.

## Delete Policy

Swagger exposes `DELETE` endpoints for domains, users, aliases, alternative domains, domain managers, relays, and tokens. Runtime validation confirmed hard delete behavior for domains, users, and aliases: reads after delete return `404`. If CETIC Group prefers safer production behavior later, the provider can add a provider-level option such as `delete_strategy = "delete"` or `delete_strategy = "disable"` after validating Mailu API support.

Any delete behavior requires architecture, QA, and security approval before implementation.

## Deferred Or Unsupported

- `mailu_fetchmail`: no Swagger endpoint found.
- `mailu_server_info`: no Swagger endpoint found.
- `mailu_dkim`: no read endpoint found; DNS DKIM value is exposed on `DomainGet.dns_dkim`, and key generation is exposed as `POST /domain/{domain}/dkim`.
