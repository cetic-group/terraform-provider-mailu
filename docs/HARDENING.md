# Hardening Guide

This guide records the production hardening rules for the Mailu Terraform provider.

## Current Hardening Controls

- API errors are redacted before they are returned as Terraform diagnostics.
- API tokens, raw passwords, auto-reply bodies, password hashes, and SMTP credentials in URLs are redacted from error messages.
- `401`, `403`, `400`, and `404` API responses are not retried.
- `429` and `5xx` API responses are retried with bounded attempts only for idempotent methods.
- `Retry-After` is parsed for rate limit responses.
- Import IDs are validated before they are written to Terraform state.
- Creates store the known resource identity if Mailu accepts the create request but the immediate read-back fails.
- `mailu_token.token` is not persisted in Terraform state after creation.
- `mailu_relay.smtp` rejects URLs containing embedded credentials.

## Production Rollout Checklist

- Use a remote Terraform backend with encryption at rest and strict access control.
- Do not commit `.tfvars`, `.env`, state files, plans, or provider release artifacts with secrets.
- Run `terraform plan` and review all deletes before production apply.
- Use disposable domains for acceptance tests.
- Import existing production Mailu objects before managing them with Terraform.
- Verify that `MAILU_ACC_DOMAIN` is not a production domain when running acceptance tests.
- Keep `insecure_skip_tls_verify` disabled outside lab environments.
- Rotate Mailu API tokens used during testing.
- Review CI logs and release artifacts for accidental secret exposure before publishing a release.

## Known Limitations

- Terraform state stores sensitive values in clear text unless the backend protects them. The provider avoids persisting generated Mailu token values for this reason.
- DNS records are not managed by this provider. Use DNS providers and the `mailu_dkim` data source or domain DNS outputs.
- `mailu_fetchmail` and `mailu_server_info` are not implemented because Mailu exposes no API endpoints for them in this installation.
- Mailu object identity is normalized to lowercase for domains and email addresses.
- Dependencies between domains, users, aliases, managers, and tokens are represented through Terraform references in configuration, not enforced by the Mailu provider schema.
- Release artifacts are checksummed, but signing and provenance attestations are deferred until CETIC Group defines a signing key policy.
- GitHub Actions are not yet pinned by commit SHA; this must be completed before a public Terraform Registry release.
