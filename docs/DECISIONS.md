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

Partially resolved: Swagger confirms an `Authorization` header is mandatory.

Decision still needed: confirm whether the value must be `Bearer <token>` or the raw token with a development token.

### Delete Strategy

Partially resolved: Swagger exposes `DELETE` endpoints.

Decision still needed: validate runtime behavior and choose between hard delete, disable-on-delete, or provider-level configurable strategy for production resources.

### DNS Ownership

Proposed direction: DNS records stay outside this provider. Mailu DNS values should be exposed as computed fields/data sources for DNS providers to consume.

Decision still needed: confirm this after `mailu_dkim`/domain DNS data source design.

### Password Updates

Partially resolved: Mailu uses `raw_password` for writes and returns `UserGet.password` as a hash.

Decision still needed: define Terraform state behavior for password changes without exposing the hash or causing perpetual diffs.
