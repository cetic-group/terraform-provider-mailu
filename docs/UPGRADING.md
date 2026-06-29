# Upgrade Guide

This guide records upgrade notes for released versions.

## From pre-release to 0.1.0

`0.1.0` is the first planned internal release candidate.

Expected upgrade work:

- Review provider configuration and ensure `endpoint` and `token` come from variables or environment variables.
- Re-run `terraform init` after installing release artifacts from the internal mirror or local plugin path.
- Review plans before applying because Mailu deletes are hard deletes.
- Keep `raw_password`, provider `token`, generated `mailu_token.token`, and Terraform state protected.

## Breaking Changes

No released breaking changes yet.

Future breaking changes must include:

- A migration note.
- Impacted resources/data sources.
- Import/state handling instructions.
- Rollback guidance where practical.
