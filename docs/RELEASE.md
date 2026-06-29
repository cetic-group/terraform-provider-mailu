# Release Process

This project uses semantic versioning and GoReleaser.

Release artifacts are published through GitHub Releases and, for stable public versions, the Terraform Registry. Private provider mirrors may be used for pre-release validation.

## Versioning

Use tags in the format:

```text
vMAJOR.MINOR.PATCH
```

Rules:

- Patch: bug fixes with no schema or behavior break.
- Minor: new resources, data sources, attributes, or compatible behavior.
- Major: breaking schema, state, import, or behavior changes.

## Local Release Check

Run before tagging:

```shell
make fmt-check
make test
make build
make docs
git diff --exit-code docs
make secret-scan
make release-snapshot
```

Review generated artifacts under `dist/`.

## Tag Release

```shell
git tag v0.1.0-rc.1
git push origin v0.1.0-rc.1
```

The release workflow builds platform binaries, publishes SHA256 checksums, signs the checksum file with GPG, publishes the Terraform Registry manifest, and creates provenance attestations.

## Artifacts

Each release must include:

- Provider zip archives for supported platforms.
- SHA256 checksum file.
- GPG detached signature for the SHA256 checksum file.
- Terraform Registry manifest.
- Provenance attestations.
- Generated changelog/release notes.

## Registry Or Mirror

Until Terraform Registry publication is complete, install from a provider mirror or local plugin path.

The Terraform provider source remains:

```text
cetic-group/mailu
```

This keeps future Terraform Registry publication compatible with existing configurations.

Private GitHub Releases are not a Terraform provider registry. Terraform cannot automatically download the provider from a private GitHub release with only the `source` address. Use one of these distribution paths until the public Registry listing is live:

- Local plugin installation from the release archive.
- A filesystem or network provider mirror.
- A future private registry or public Terraform Registry publication after approval.

See [Private Provider Installation](PRIVATE_INSTALL.md) for the local installation procedure.

See [Public Publication And GPG Signing](PUBLICATION.md) for Terraform Registry publication, GPG key creation, Vault storage, rotation, and revocation.

## Release Candidate Validation

For `v0.1.0-rc.1`, confirm:

- The GitHub Actions release workflow completed successfully.
- Release assets exist for Linux, macOS, and Windows on `amd64` and `arm64`.
- The `SHA256SUMS` asset is attached.
- The `SHA256SUMS.sig` asset is attached for public releases.
- The Terraform Registry manifest asset is attached for public releases.
- The release is marked as a pre-release.
- Documentation and changelog match the implemented provider surface.

## Security Review

Before publishing:

- Confirm no `.env`, `.tfvars`, token, password, or generated token value is committed.
- Confirm release logs do not include acceptance-test secrets.
- Confirm `mailu_token.token`, `mailu_user.raw_password`, provider `token`, and `reply_body` remain sensitive.
- Confirm checksums are published with the release.
- Confirm checksum signatures are published with public releases.
- Confirm the GPG fingerprint matches the approved project release key.
