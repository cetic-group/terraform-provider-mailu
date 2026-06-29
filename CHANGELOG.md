# Changelog

All notable changes to this project are documented in this file.

This project follows semantic versioning.

## [Unreleased]

### Added

- Public Terraform Registry publication workflow hardening.
- GPG signing configuration for release checksums.
- Terraform Registry manifest release asset.
- GitHub artifact provenance attestations.
- Secret scanning target and CI job.
- Public publication, GPG key, Vault storage, rotation, and revocation documentation.

## [0.1.0-rc.1] - 2026-06-29

Initial release candidate for Mailu automation.

### Added

- Provider configuration, Mailu API client, diagnostics, redaction, retries, and local install workflow.
- Resources: `mailu_domain`, `mailu_user`, `mailu_alias`, `mailu_alternative_domain`, `mailu_domain_manager`, `mailu_relay`, and `mailu_token`.
- Data sources: `mailu_domain`, `mailu_user`, and `mailu_dkim`.
- DNS integration guidance and Terraform examples.
- CI and release engineering configuration.
