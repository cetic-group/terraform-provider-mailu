# Decisions

This file records project decisions that should remain visible without expanding the roadmap.

## Accepted

### Provider Framework

Use `terraform-plugin-framework` for the provider implementation.

Reason: it is HashiCorp's recommended framework for new Terraform providers and supports modern provider protocol features.

### Provider Identity

Use:

- Terraform source: `cetic-group/mailu`
- Provider type: `mailu`
- Go module: `github.com/cetic-group/terraform-provider-mailu`

### API First

Do not implement Terraform resources from assumptions. Capture the Mailu API contract first in [API.md](API.md).

Reason: Terraform resource behavior must be idempotent and importable. That requires confirmed Mailu read, update, delete, and error semantics.

### Secrets

Treat API tokens, passwords, private keys, and generated secrets as sensitive.

Reason: Terraform state, logs, diagnostics, and CI artifacts can otherwise expose production mail infrastructure.

### API Base Path

Use `https://mail.cetic-group.com/api/v1` as the Mailu API endpoint for this installation.

Reason: `https://mail.cetic-group.com/api/v1/swagger.json` exposes Swagger 2.0 with `basePath` set to `/api/v1`.

### MVP Resources

Keep the MVP focused on:

- `mailu_domain`
- `mailu_user`
- `mailu_alias`

Reason: Swagger confirms full list/read/create/update/delete endpoints for these objects.

### Deferred Resources

Defer or block:

- `mailu_fetchmail`: no endpoint exposed in Swagger.
- `mailu_server_info`: no endpoint exposed in Swagger.
- standalone `mailu_dkim`: no read endpoint exposed; DKIM DNS data is available through domain reads and key generation is an action.

## Open

### Authentication Header

Resolved for MVP implementation: use `Authorization: Bearer <token>`.

Reason: Swagger confirms an `Authorization` header with a `Bearer` security scheme. Runtime validation on 2026-06-29 showed both raw token and `Bearer <token>` return `200`, while invalid tokens return `403`.

### Delete Strategy

Resolved for MVP implementation: Terraform delete maps to Mailu `DELETE` for MVP resources.

Reason: runtime validation on 2026-06-29 showed `DELETE` on domains, users, and aliases returns `200`, and subsequent reads return `404`.

Production note: because this is hard delete behavior, acceptance tests must use disposable objects and production applies require plan review.

### DNS Ownership

Proposed direction: DNS records stay outside this provider. Mailu DNS values should be exposed as computed fields/data sources for DNS providers to consume.

Decision still needed: confirm this after `mailu_dkim`/domain DNS data source design.

### Password Updates

Resolved for MVP implementation: expose `raw_password` as sensitive write input and do not expose `UserGet.password`.

Reason: Mailu uses `raw_password` for writes and returns `UserGet.password` as a string hash. The provider must not compare the raw password with the returned hash and must redact the hash from diagnostics/logs.
