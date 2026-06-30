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
- `mailu_token.token` is stored in Terraform state as a sensitive value (Mailu returns it only at creation time); protect the state backend accordingly.
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

- Terraform state stores sensitive values in clear text unless the backend protects them. Generated Mailu token values are stored in state as sensitive values because Mailu only returns them at creation time; protect the state backend accordingly.
- DNS records are not managed by this provider. Use DNS providers and the `mailu_dkim` data source or domain DNS outputs.
- `mailu_fetchmail` and `mailu_server_info` are not implemented because Mailu exposes no API endpoints for them in this installation.
- Mailu object identity is normalized to lowercase for domains and email addresses.
- Dependencies between domains, users, aliases, managers, and tokens are represented through Terraform references in configuration, not enforced by the Mailu provider schema.
- Release artifacts are checksummed and public release checksums are signed with the project GPG release key.
- GitHub Actions are pinned by commit SHA for CI and release workflows.
- GPG signing material is governed by [Public Publication And GPG Signing](PUBLICATION.md) and stored in Vault as the source of authority.
