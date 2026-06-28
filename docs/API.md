# Mailu API Contract

Discovery source: `https://mail.cetic-group.com/api/v1/swagger.json`

Discovery date: 2026-06-28

Spec format: Swagger 2.0

API title/version: `Mailu API` / `1.0`

Base path: `/api/v1`

Status: partially confirmed from the public Swagger document. Runtime behavior with a valid development token still needs validation before CRUD implementation.

## Authentication

Confirmed by Swagger:

- Security scheme: `Bearer`
- Type: `apiKey`
- Header name: `Authorization`
- All operations require the `Authorization` header.

Confirmed by unauthenticated request:

```json
{"message": "A valid Authorization header is mandatory"}
```

Still to validate with a development token:

- Whether the header value must be `Bearer <token>` or the raw token.
- Whether invalid tokens return `403` consistently.
- Whether domain managers and global admins have different API permissions.

## Content Type

Confirmed by Swagger:

- Consumes: `application/json`
- Produces: `application/json`

## Error Model

Most validation and not-found errors reference:

```json
{
  "code": 400,
  "message": "error message"
}
```

Common status codes:

- `200`: success
- `400`: input validation exception
- `401`: authorization header missing
- `403`: invalid authorization header
- `404`: object not found
- `409`: duplicate object

## Endpoint Mapping

| Terraform object | API operations | Endpoint(s) | Status |
| --- | --- | --- | --- |
| `mailu_domain` | list, read, create, update, delete | `GET/POST /domain`, `GET/PATCH/DELETE /domain/{domain}` | confirmed by Swagger |
| `mailu_user` | list, read, create, update, delete | `GET/POST /user`, `GET/PATCH/DELETE /user/{email}`, `GET /domain/{domain}/users` | confirmed by Swagger |
| `mailu_alias` | list, read, create, update, delete | `GET/POST /alias`, `GET/PATCH/DELETE /alias/{alias}`, `GET /alias/destination/{domain}` | confirmed by Swagger |
| `mailu_alternative_domain` | list, read, create, delete | `GET/POST /alternative`, `GET/DELETE /alternative/{alt}` | confirmed by Swagger |
| `mailu_domain_manager` | list, read, create, delete | `GET/POST /domain/{domain}/manager`, `GET/DELETE /domain/{domain}/manager/{email}` | confirmed by Swagger |
| `mailu_relay` | list, read, create, update, delete | `GET/POST /relay`, `GET/PATCH/DELETE /relay/{name}` | confirmed by Swagger |
| `mailu_token` | list, read, create, update, delete | `GET/POST /token`, `GET/PATCH/DELETE /token/{token_id}`, `GET/POST /tokenuser/{email}` | confirmed by Swagger |
| `mailu_dkim` | generate only; DNS values readable from domain | `POST /domain/{domain}/dkim`, `GET /domain/{domain}` | partial |
| `mailu_forward` | user fields only | `GET/PATCH /user/{email}` | model as part of `mailu_user` first |
| `mailu_fetchmail` | none | none | blocked; no Swagger endpoint |
| `mailu_server_info` | none | none | blocked; no Swagger endpoint |

## Schemas

### Domain

Create schema: `Domain`

Required:

- `name`

Optional:

- `comment`
- `max_users`
- `max_aliases`
- `max_quota_bytes`
- `signup_enabled`
- `alternatives`

Read schema: `DomainGet`

Read-only or computed fields:

- `managers`
- `dns_autoconfig`
- `dns_mx`
- `dns_spf`
- `dns_dkim`
- `dns_dmarc`
- `dns_dmarc_report`
- `dns_tlsa`

Update schema: `DomainUpdate`

Mutable fields:

- `comment`
- `max_users`
- `max_aliases`
- `max_quota_bytes`
- `signup_enabled`
- `alternatives`

Notes:

- Domain rename is not exposed; changing `name` should force replacement.
- DKIM generation is an action, not a read resource. DNS DKIM output is exposed through `DomainGet.dns_dkim`.

### User

Create schema: `UserCreate`

Required:

- `email`
- `raw_password`

Optional:

- `comment`
- `quota_bytes`
- `global_admin`
- `enabled`
- `change_pw_next_login`
- `enable_imap`
- `enable_pop`
- `allow_spoofing`
- `forward_enabled`
- `forward_destination`
- `forward_keep`
- `reply_enabled`
- `reply_subject`
- `reply_body`
- `reply_startdate`
- `reply_enddate`
- `displayed_name`
- `spam_enabled`
- `spam_mark_as_read`
- `spam_threshold`

Read schema: `UserGet`

Additional read fields:

- `password`
- `quota_bytes_used`

Update schema: `UserUpdate`

Mutable fields are the create optional fields plus `raw_password`.

Notes:

- `raw_password` is write-only input and must be sensitive.
- `UserGet.password` is a password hash and must not be exposed as a normal Terraform attribute.
- User ID/import ID should be the full email address.

### Alias

Create schema: `Alias`

Required:

- `email`

Optional:

- `destination`
- `comment`
- `wildcard`

Update schema: `AliasUpdate`

Mutable fields:

- `destination`
- `comment`
- `wildcard`

Notes:

- Alias ID/import ID should be the full alias email address.
- Terraform should expose `destination` as a set/list of email addresses.

### Alternative Domain

Create schema: `AlternativeDomain`

Required:

- `name`
- `domain`

Notes:

- No update endpoint is exposed.
- Import ID can be the alternative domain name.

### Domain Manager

Create schema: `ManagerCreate`

Required:

- `user_email`

Path input:

- `domain`

Notes:

- No update endpoint is exposed.
- Import ID can be `<domain>/<email>`.

### Relay

Create schema: `Relay`

Required:

- `name`

Optional:

- `smtp`
- `comment`

Update schema: `RelayUpdate`

Mutable fields:

- `smtp`
- `comment`

### Token

Create schema: `TokenPost` or `TokenPost2`

Required:

- `email` for `POST /token`

Optional:

- `comment`
- `AuthorizedIP`

Read schema: `TokenGetResponse`

Read fields:

- `id`
- `email`
- `comment`
- `AuthorizedIP`
- `Created`
- `Last edit`

Create response additionally returns:

- `token`

Security note:

- `token` is returned only on create and must be sensitive/write-only in Terraform state handling.

## Pagination

No pagination parameters are present in the Swagger paths.

Runtime validation still needs to confirm behavior for large collections.

## Import IDs

Recommended import IDs from Swagger path identifiers:

- `mailu_domain`: `example.com`
- `mailu_user`: `user@example.com`
- `mailu_alias`: `alias@example.com`
- `mailu_alternative_domain`: `alias-domain.example`
- `mailu_domain_manager`: `example.com/user@example.com`
- `mailu_relay`: `example.com`
- `mailu_token`: token record `id`

## Redaction Rules

When adding captured API examples:

- Redact API tokens.
- Redact `raw_password`.
- Redact token create response values.
- Redact password hashes from `UserGet.password`.
- Redact private keys.
- Keep domains and email addresses only when they are intentional examples.

## Runtime Validation Still Required

Before implementation:

- Test authenticated list/read calls with a development token.
- Test whether `Authorization: Bearer <token>` or `Authorization: <token>` is required.
- Test actual create/update/delete against a disposable domain.
- Capture real success payloads for write operations.
- Capture real validation failures.
- Confirm whether `PATCH` behaves as partial update.
- Confirm whether `DELETE` is hard delete.
- Confirm whether collection responses are arrays for all list endpoints.

## Agent Review Findings

Senior Developer Architect:

- MVP resources are feasible from Swagger: `mailu_domain`, `mailu_user`, `mailu_alias`.
- `mailu_fetchmail` and `mailu_server_info` must be deferred or removed because no endpoints are exposed.
- `mailu_dkim` should start as a domain DNS data source field, not a standalone CRUD resource.

Senior QA:

- Swagger discovery is reproducible with `curl https://mail.cetic-group.com/api/v1/swagger.json`.
- Acceptance testing still needs a dedicated token and disposable domain before CRUD behavior can be validated.
- Tests must include unauthenticated, invalid-token, duplicate, not-found, and validation-error cases.

Senior Application Security:

- `Authorization`, `raw_password`, token create responses, and user password hashes are sensitive.
- Token management is high risk because create responses expose newly generated tokens.
- Destructive operations should not be considered production-safe until runtime delete behavior is validated.
