# Release Process

This project uses semantic versioning and GoReleaser.

Initial channel decision: publish release artifacts to the CETIC Group internal mirror or GitHub Releases first. Public Terraform Registry publication is deferred until CETIC Group explicitly approves public distribution.

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
make release-snapshot
```

Review generated artifacts under `dist/`.

## Tag Release

```shell
git tag v0.1.0
git push origin v0.1.0
```

The release workflow builds platform binaries and publishes SHA256 checksums.

## Artifacts

Each release must include:

- Provider zip archives for supported platforms.
- SHA256 checksum file.
- Generated changelog/release notes.

## Registry Or Mirror

Until public publication is approved, install from the internal mirror or local plugin path.

The Terraform provider source remains:

```text
cetic-group/mailu
```

This keeps future Terraform Registry publication compatible with existing configurations.

## Security Review

Before publishing:

- Confirm no `.env`, `.tfvars`, token, password, or generated token value is committed.
- Confirm release logs do not include acceptance-test secrets.
- Confirm `mailu_token.token`, `mailu_user.raw_password`, provider `token`, and `reply_body` remain sensitive.
- Confirm checksums are published with the release.
