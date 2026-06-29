# Public Publication And GPG Signing

This document describes the CETIC Group process for publishing `cetic-group/mailu` to the public Terraform Registry.

Public publication requires explicit CETIC Group approval. Do not publish a public release until the GitHub repository, GPG key, Terraform Registry namespace, and release workflow have been reviewed.

## Publication Prerequisites

- GitHub repository is public and named `terraform-provider-mailu`.
- Terraform Registry namespace `cetic-group` is available to CETIC Group.
- Provider source remains `cetic-group/mailu`.
- GPG signing key is created, approved, stored in Vault, and added to Terraform Registry.
- GitHub release workflow has access to signing material through approved secrets.
- Release artifacts include platform ZIP archives, `SHA256SUMS`, `SHA256SUMS.sig`, the Registry manifest, and provenance attestations.

## GPG Key Purpose

GPG signing proves that release checksums were produced with a CETIC-controlled signing key.

Checksums alone prove that an archive matches a published checksum. A detached GPG signature also proves that the checksum file was signed by the holder of the CETIC release key. Terraform Registry validates provider releases against the public GPG key registered for the provider.

## Create The GPG Key

Create a dedicated release key. Do not use a personal developer key.

Recommended identity:

```text
CETIC Group Terraform Provider Releases <release@cetic-group.com>
```

Recommended parameters:

- RSA 4096.
- Expiration: 12 to 24 months.
- Strong passphrase stored in Vault.
- Two-person review of the key fingerprint before Registry publication.

Example:

```shell
gpg --full-generate-key
```

List the key and record the full fingerprint:

```shell
gpg --list-secret-keys --keyid-format LONG
gpg --fingerprint release@cetic-group.com
```

Export the public key:

```shell
gpg --armor --export release@cetic-group.com > cetic-terraform-provider-public.gpg
```

Export the private key:

```shell
gpg --armor --export-secret-keys release@cetic-group.com > cetic-terraform-provider-private.gpg
```

Generate and store a revocation certificate immediately after key creation:

```shell
gpg --output cetic-terraform-provider-revocation.asc --gen-revoke release@cetic-group.com
```

## Store The Key In Vault

Vault is the source of authority for signing material. GitHub Secrets are an operational copy only when Vault OIDC integration is not yet available.

Recommended Vault paths:

```text
kv/release/terraform-provider-mailu/gpg/private_key
kv/release/terraform-provider-mailu/gpg/passphrase
kv/release/terraform-provider-mailu/gpg/public_key
kv/release/terraform-provider-mailu/gpg/fingerprint
kv/release/terraform-provider-mailu/gpg/revocation_certificate
```

Vault requirements:

- KV v2 versioning enabled.
- Audit devices enabled.
- Read access to private key and passphrase limited to the release process.
- Human break-glass access documented and reviewed.
- Rotation and revocation procedure approved by CETIC Group security.

## Configure GitHub

Use a protected GitHub environment named `release` with required reviewers.

Add these environment secrets:

```text
GPG_PRIVATE_KEY
GPG_PASSPHRASE
GPG_FINGERPRINT
```

Values:

- `GPG_PRIVATE_KEY`: full ASCII-armored private key.
- `GPG_PASSPHRASE`: passphrase for the key.
- `GPG_FINGERPRINT`: full approved key fingerprint.

The private key and passphrase must come from Vault. Do not paste them into issues, logs, artifacts, `.env` files, or repository files.

## Configure Terraform Registry

In Terraform Registry:

1. Sign in with the GitHub account that can access the CETIC Group repository.
2. Confirm the `cetic-group` namespace ownership.
3. Publish a provider from the public GitHub repository `terraform-provider-mailu`.
4. Add the GPG public key exported from Vault.
5. Verify the fingerprint matches `GPG_FINGERPRINT`.
6. Confirm provider source is `cetic-group/mailu`.

Terraform Registry uses GitHub releases and the registered GPG key to ingest provider versions.

## Release Procedure

Before tagging:

```shell
make fmt-check
make test
make build
make docs
git diff --exit-code docs
make secret-scan
```

Then tag a stable release:

```shell
git tag v0.1.0
git push origin v0.1.0
```

The release workflow imports the GPG key, verifies the expected fingerprint, runs tests and documentation checks, scans committed history for secrets, builds platform archives, creates SHA256 checksums, signs the checksum file, publishes the Terraform Registry manifest, and creates GitHub provenance attestations.

## Verify A Release

Download these release assets:

```text
terraform-provider-mailu_<version>_SHA256SUMS
terraform-provider-mailu_<version>_SHA256SUMS.sig
cetic-terraform-provider-public.gpg
```

Verify the signature:

```shell
gpg --import cetic-terraform-provider-public.gpg
gpg --verify terraform-provider-mailu_<version>_SHA256SUMS.sig terraform-provider-mailu_<version>_SHA256SUMS
```

Verify an archive checksum:

```shell
shasum -a 256 terraform-provider-mailu_<version>_darwin_arm64.zip
grep terraform-provider-mailu_<version>_darwin_arm64.zip terraform-provider-mailu_<version>_SHA256SUMS
```

## Rotation

Rotate the GPG key every 12 to 24 months, before expiration, immediately after suspected exposure, or after release team ownership changes if access cannot be proven clean.

Rotation steps:

1. Create a new key.
2. Store private key, passphrase, public key, fingerprint, and revocation certificate in Vault.
3. Add the new public key to Terraform Registry.
4. Update GitHub release environment secrets from Vault.
5. Create a test release candidate signed by the new key.
6. Verify `SHA256SUMS.sig` with the new public key.
7. Record the rotation in release notes.

Keep old public keys available while releases signed with them remain supported.

## Revocation

Trigger revocation if the private key or passphrase may be exposed.

Emergency steps:

1. Disable the GitHub release workflow or remove release environment approvals.
2. Revoke GitHub environment secrets.
3. Disable Vault CI access to the compromised key.
4. Publish the revocation certificate.
5. Remove or mark compromised releases according to GitHub and Terraform Registry guidance.
6. Create a new key and update Terraform Registry.
7. Publish a new fixed release.
8. Publish a security advisory if public users may be affected.

Do not overwrite an already published provider version unless CETIC Group security explicitly approves the incident response. Prefer a new release version.
