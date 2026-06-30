# Changelog

All notable changes to this project are documented in this file.

This project follows semantic versioning.

## [Unreleased]

### Added

- `mailu_domains`, `mailu_users`, and `mailu_aliases` list data sources for inventorying existing objects when generating Terraform import blocks.
- Plan-time schema validators for email and domain-name formats, token authorized IP/CIDR values, and the user spam threshold range.
- Acceptance tests for `mailu_domain`, `mailu_user`, and `mailu_alias` (create, update/drift, import, destroy) and an opt-in `acceptance` GitHub workflow.
- Public Terraform Registry publication workflow hardening.
- GPG signing configuration for release checksums.
- Terraform Registry manifest release asset.
- GitHub artifact provenance attestations.
- Secret scanning target and CI job.
- Public publication, GPG key, Vault storage, rotation, and revocation documentation.

### Changed

- All `Optional`+`Computed` resource attributes now use the `UseStateForUnknown` plan modifier, eliminating spurious "(known after apply)" diffs on unrelated changes.
- `mailu_token.token` is now stored in Terraform state as a sensitive value (Mailu returns it only at creation time) instead of being discarded, so it can be consumed by outputs and other resources. Protect the state backend accordingly.

## [0.1.0-rc.1] - 2026-06-29

Initial release candidate for Mailu automation.

### Added

- Provider configuration, Mailu API client, diagnostics, redaction, retries, and local install workflow.
- Resources: `mailu_domain`, `mailu_user`, `mailu_alias`, `mailu_alternative_domain`, `mailu_domain_manager`, `mailu_relay`, and `mailu_token`.
- Data sources: `mailu_domain`, `mailu_user`, and `mailu_dkim`.
- DNS integration guidance and Terraform examples.
- CI and release engineering configuration.
