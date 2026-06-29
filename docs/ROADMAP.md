# Roadmap

This roadmap tracks implementation of the `cetic-group/mailu` Terraform provider.

Process rule: phases are handled in order. A phase is complete only when its exit criteria are satisfied and architecture, QA, and security reviews have no blocking findings.

## Phase 1 - Project Baseline

Status: complete.

Goal: establish the provider project, documentation, review process, and local development prerequisites.

Tasks:

- Create the provider repository scaffold.
- Use Go module `github.com/cetic-group/terraform-provider-mailu`.
- Use provider source `registry.terraform.io/cetic-group/mailu`.
- Configure provider type name `mailu`.
- Add provider configuration through HCL and environment variables.
- Add the shared Mailu HTTP client scaffold.
- Add initial examples, Terraform provider documentation, and project documentation.
- Add project agents for architecture, QA, and application security.
- Install Go 1.25.8+ locally.
- Confirm Terraform 1.8+ is available.
- Define local test variables with `MAILU_ENDPOINT`, `MAILU_API_TOKEN`, and `MAILU_ACC_DOMAIN`.
- Store local secrets outside Git.

Exit criteria:

- README links to project, Terraform, contribution, decision, API, and roadmap docs.
- Terraform provider documentation exists under `docs/index.md`.
- Planned resource and data source pages exist under `docs/resources` and `docs/data-sources`.
- `.env.local` and other local secrets are ignored by Git.
- `go version` reports Go 1.25.8+.
- `terraform version` reports Terraform 1.8+.
- `go mod tidy`, `gofmt -w .`, and `go test ./...` pass.
- Provider can be installed locally.
- Terraform can load the provider with `source = "cetic-group/mailu"`.

## Phase 2 - Mailu API Discovery

Status: complete for MVP resources; extended resource runtime validation deferred to their implementation phases.

Goal: capture the real Mailu API contract before implementing Terraform resources.

Tasks:

- Confirm the exact authentication mechanism.
- Confirm endpoint paths for domains, users, aliases, DKIM, forwards, fetchmail, and server metadata.
- Capture list, read, create, update, and delete behavior.
- Capture response payload shape for successful and failed requests.
- Capture pagination behavior, if any.
- Capture stable identifiers and natural import IDs.
- Confirm whether each object can be read after creation.
- Confirm whether updates are patch, replace, or delete/recreate operations.
- Confirm delete behavior.
- Confirm password handling and whether password hashes are returned.
- Confirm exposed fields for admin flags, quota, enabled state, display name, spam threshold, and destinations.
- Confirm Mailu version and API compatibility assumptions.
- Write redacted findings in `docs/API.md`.
- Update `docs/RESOURCE_MODEL.md` with confirmed and unsupported fields.
- Update `docs/DECISIONS.md` with resolved authentication, delete, DNS, and password decisions.
- Run architecture, QA, and security reviews.

Exit criteria:

- Every MVP resource has endpoint mapping.
- Import IDs are specified.
- API gaps are documented with a fallback decision.
- Unsupported resources are explicitly marked as blocked or deferred.
- No MVP implementation proceeds from unverified endpoint assumptions.

Current result:

- Swagger discovery is complete.
- Runtime auth and MVP CRUD validation completed on 2026-06-29.
- `mailu_domain`, `mailu_user`, and `mailu_alias` were validated with temporary objects.
- Extended resources such as relay, token, domain manager, alternative domain, and DKIM generation remain mapped from Swagger but need dedicated runtime validation before implementation.

## Phase 3 - MVP Design Freeze

Status: complete.

Goal: freeze the first implementable provider surface before coding resources.

Tasks:

- Finalize schemas for `mailu_domain`, `mailu_user`, and `mailu_alias`.
- Mark each schema attribute as required, optional, computed, sensitive, or replacement-forcing.
- Define Terraform IDs and import ID formats.
- Define normalization rules for domains and email addresses.
- Define password update behavior.
- Define delete behavior.
- Define drift behavior for every MVP resource.
- Define acceptance test fixtures and cleanup rules.
- Update Terraform documentation pages from confirmed schemas.

Exit criteria:

- `docs/RESOURCE_MODEL.md` matches confirmed API behavior.
- `docs/resources/*.md` matches MVP schemas.
- `docs/data-sources/*.md` matches MVP schemas.
- `docs/DECISIONS.md` has no open blocker for MVP implementation.
- Architecture, QA, and security reviews approve the MVP design.

## Phase 4 - Provider Foundation

Status: complete.

Goal: make the provider framework, client, diagnostics, and test harness implementation-ready.

Tasks:

- Add client tests with `httptest`.
- Implement API error types.
- Add retry and timeout configuration.
- Add Terraform diagnostics helpers.
- Add redaction helpers for diagnostics and logs.
- Add user agent with provider version.
- Add environment-variable based test configuration.
- Add acceptance test harness gated by `TF_ACC=1`.
- Add Makefile targets for `fmt`, `test`, and acceptance tests.
- Add local provider install instructions for the current platform.
- Add provider-level `insecure_skip_tls_verify` only if needed for lab environments.

Acceptance test variables:

- `MAILU_ENDPOINT`
- `MAILU_API_TOKEN`
- `MAILU_ACC_DOMAIN`

Exit criteria:

- Unit tests cover client success and error responses.
- Unit tests cover authentication header behavior.
- Unit tests cover URL resolution.
- Unit tests cover redaction behavior.
- Provider configuration tests cover missing config and environment fallback.
- Acceptance tests can be run explicitly with `TF_ACC=1`.
- Acceptance tests cannot run without an explicit disposable domain.

## Phase 5 - MVP Resources

Status: complete.

Goal: implement the first useful provider surface.

Resources:

- `mailu_domain`
- `mailu_user`
- `mailu_alias`

Data sources:

- `mailu_domain`
- `mailu_user`

Expected examples:

```hcl
resource "mailu_domain" "cetic" {
  name = "example.com"
}

resource "mailu_user" "admin" {
  email        = "admin@example.com"
  raw_password = var.admin_password
  quota_bytes  = 1073741824
  enabled      = true
  global_admin = true
}

resource "mailu_alias" "postmaster" {
  email       = "postmaster@example.com"
  destination = ["admin@example.com"]
}
```

Import IDs:

- Domain: `example.com`
- User: `admin@example.com`
- Alias: `postmaster@example.com`

Exit criteria:

- Existing Mailu objects can be imported.
- Terraform detects drift.
- Passwords remain sensitive and are not read back from the API.
- Create, read, update, delete, and import are tested for every MVP resource.
- Examples run against the test environment.
- Terraform documentation no longer marks implemented MVP resources as planned.
- Architecture, QA, and security reviews are complete.

Current result:

- `mailu_domain`, `mailu_user`, and `mailu_alias` are implemented with CRUD and import support.
- `mailu_domain` and `mailu_user` data sources are implemented.
- Runtime validation on 2026-06-29 created, read, imported, detected drift on, and destroyed temporary MVP objects under `tf-phase5-*.<MAILU_ACC_DOMAIN>`.
- Terraform plan output keeps `raw_password` and `reply_body` sensitive.
- Documentation and examples now describe the implemented MVP surface.
- Architecture, QA, and application security reviews found no blocking issue for the MVP scope.

## Phase 6 - Extended Mail Resources

Status: complete.

Goal: expand provider coverage after the MVP is stable.

Candidate resources:

- `mailu_forward`, modeled first as part of `mailu_user`.
- `mailu_relay`.
- `mailu_token`, with strict sensitive handling.
- `mailu_alternative_domain`.
- `mailu_domain_manager`.
- `mailu_fetchmail`, blocked because no API endpoint is exposed.

Candidate data sources:

- `mailu_dkim`, backed by domain DNS fields.
- `mailu_server_info`, blocked because no API endpoint is exposed.

Exit criteria:

- Each added object follows the same discovery, design, test, documentation, and review process as MVP resources.
- Runtime validation exists before implementation for every extended resource.
- Token resources do not expose generated token values outside sensitive state handling.

Current result:

- `mailu_alternative_domain`, `mailu_domain_manager`, `mailu_relay`, and `mailu_token` are implemented.
- `mailu_dkim` is implemented as a read-only data source backed by domain DNS fields.
- User forwarding remains modeled through `mailu_user` fields.
- `mailu_fetchmail` and `mailu_server_info` remain blocked because the API exposes no endpoints for them.
- Runtime validation on 2026-06-29 confirmed create/read/update/delete behavior for implemented extended resources with temporary objects.
- DKIM generation was runtime-validated as an API action, but no Terraform resource was added because key rotation is not a CRUD lifecycle.
- Token generated values are marked sensitive and only available from create/import-safe state handling.
- Terraform apply, stable plan, import, and destroy pass for the full phase 6 surface.
- A Terraform apply uncovered that Mailu returns numeric token IDs despite Swagger declaring strings; the client now accepts both string and numeric token IDs.

## Phase 7 - DNS Integration Patterns

Status: complete.

Goal: document how this provider composes with DNS providers.

Tasks:

- Produce examples for MX, SPF, DKIM, DMARC, MTA-STS, and autoconfig records.
- Add an IONOS-oriented example for deployments using IONOS DNS automation.
- Document ownership boundaries: Mailu provider manages Mailu state, DNS provider manages DNS state.
- Add an example that consumes `mailu_dkim` or domain DNS output in DNS records, if implemented.
- Document which DNS records are mandatory, recommended, or optional.

Exit criteria:

- A new domain can be onboarded with Mailu and DNS using Terraform examples.
- DNS examples do not require secrets in source files.

Current result:

- DNS ownership boundaries are documented in `docs/DNS.md` and `docs/DECISIONS.md`.
- Generic DNS composition example exists under `examples/dns/generic`.
- IONOS-oriented payload example exists under `examples/dns/ionos`.
- DNS docs cover MX, SPF, DKIM, DMARC, MTA-STS, TLS reporting, autoconfig, TLSA, and record priority.
- Examples use provider variables and environment variables for secrets; no DNS or Mailu secret is committed.

## Phase 8 - Release Engineering

Status: complete.

Goal: make releases repeatable and installable.

Tasks:

- Add GitHub Actions CI.
- Add `goreleaser`.
- Generate provider docs with `terraform-plugin-docs`.
- Keep Terraform Registry documentation under `docs/index.md`, `docs/resources`, and `docs/data-sources`.
- Add changelog.
- Add semantic versioning.
- Decide Terraform Registry, GitHub release asset, or private mirror distribution.
- Configure release signing or checksum publication.
- Add release notes template.
- Add upgrade notes for breaking changes.
- Publish `v0.1.0` through GitHub Releases and Terraform Registry when publication gates are approved.
- Run release readiness review.

Exit criteria:

- Tagged releases publish checksums and platform binaries.
- The provider can be installed from the Terraform Registry, GitHub release assets, or a private mirror.
- Documentation generated for the release matches implemented schemas.
- Release process is repeatable from a clean checkout.

Current result:

- GitHub Actions CI workflow added for formatting, tests, build, and generated documentation checks.
- GitHub Actions release workflow added for semantic version tags.
- GoReleaser configuration added for Linux, macOS, and Windows `amd64`/`arm64` archives.
- SHA256 checksum publication configured through GoReleaser.
- Release process, release notes template, changelog, and upgrade guide added.
- GitHub Releases and Terraform Registry publication are the primary public distribution paths.
- `make docs` uses `terraform-plugin-docs` to regenerate provider documentation.
- `v0.1.0-rc.1` was published as a pre-release with platform archives and `SHA256SUMS` attached.
- Private installation from GitHub release assets is documented in `docs/PRIVATE_INSTALL.md`.
- Pre-release assets are installable through the local Terraform plugin directory or a provider mirror before a stable Registry release.

## Phase 9 - Hardening

Status: complete for pre-production hardening; production adoption remains Phase 10.

Goal: prepare the provider for production-grade use.

Tasks:

- Add rate limit handling.
- Add diff suppression for normalized email addresses.
- Finalize case normalization policy for domains and email addresses.
- Add state migration tests.
- Add import validation.
- Add partial-state handling for create/update failures.
- Add detailed acceptance tests against disposable domains.
- Review concurrency behavior.
- Tune timeouts and retries against real Mailu behavior.
- Add negative tests for authorization failures.
- Add state upgrade framework for future schema versions.
- Review CI logs and release artifacts for security exposure.

Exit criteria:

- Provider behavior is stable enough for production mailbox/domain management.
- Production rollout checklist exists.
- Known limitations are documented in Terraform provider docs.

Current result:

- HTTP retry handling covers `429` and `5xx` for idempotent methods, respects `Retry-After`, and does not retry authorization, client validation, create, or update failures.
- Import IDs are validated before being stored for domains, users, aliases, relays, tokens, and domain managers.
- Create operations store known identity state and emit a warning if Mailu accepts the create but the immediate read-back fails.
- API diagnostic redaction now covers raw tokens, bearer tokens, passwords, `raw_password`, `reply_body`, token fields, and SMTP URL credentials.
- `mailu_token.token` is no longer persisted in Terraform state after creation because Terraform state stores sensitive values in clear text.
- `mailu_relay.smtp` rejects URLs that contain embedded credentials.
- Unit tests cover rate limits, authorization non-retry behavior, import validation, relay SMTP credential rejection, generated token non-persistence, and redaction.
- Production rollout checklist and known limitations are documented in `docs/HARDENING.md`.
- Supply-chain signing, provenance attestations, and GitHub Actions SHA pinning are documented as pre-public-release hardening work.
- Full production migration and operational adoption remain in Phase 10.

## Phase 10 - Production Adoption

Status: complete for minimal adoption; remote production backend and broader rollout deferred.

Goal: migrate existing Mailu management to Terraform safely.

Tasks:

- Inventory existing Mailu domains, users, aliases, and related objects.
- Generate Terraform import blocks or import commands.
- Import existing objects into a non-production state first.
- Compare Terraform state with Mailu API reads.
- Prepare production state backend and access policy.
- Run `terraform plan` with no destructive changes.
- Review plan with architecture, QA, and security agents.
- Apply first production change on a low-risk object.
- Document rollback and manual recovery steps.

Exit criteria:

- Existing minimal Mailu objects are imported without unintended changes.
- Minimal imported state produces a no-op plan.
- Operational runbook exists.

Current result:

- Production adoption runbook added in `docs/PRODUCTION_ADOPTION.md`.
- Production scaffold added in `examples/production` with provider configuration, import block examples, and backend template.
- Import examples now cover all implemented resources.
- Minimal non-production adoption validation completed from `/private/tmp/mailu-terraform-adoption` for `mailu_domain.example` (`example.com`) and `mailu_user.admin` (`admin@example.com`).
- `terraform plan -refresh=true -detailed-exitcode` returned no changes for the minimal imported state.
- No destructive or mutating production apply was executed.
- Remote production backend, broader production import, no-op plan review in remote state, and first low-risk apply are deferred until the operator decides to run a wider production rollout.

## Phase 11 - Public Release And Supply Chain Hardening

Status: complete for public release readiness; external GitHub and Terraform Registry actions remain manual gates.

Goal: prepare the provider for public Terraform Registry publication and stronger release-chain guarantees.

Scope note: this phase is intentionally separate from Phase 10. Production adoption can proceed through controlled release assets or a private mirror, while public release requires additional supply-chain controls and maintainer approval.

Tasks:

- Confirm maintainer approval for public repository and Terraform Registry publication.
- Define signing key ownership, storage, rotation, and revocation policy.
- Pin GitHub Actions by commit SHA instead of mutable version tags.
- Pin GoReleaser to an explicit reviewed version.
- Add artifact signing for checksums and release archives.
- Add provenance attestations for release builds.
- Add secret scanning to CI and release workflows.
- Document public release rollback and revocation procedure.
- Validate Terraform Registry publication metadata, namespace ownership, and provider documentation.
- Run architecture, QA, and application security release reviews before first public release.

Exit criteria:

- Release workflow uses pinned, reviewed dependencies.
- Release artifacts are checksummed, signed, and accompanied by provenance attestations.
- Project signing policy is documented and approved.
- Terraform Registry publication procedure is documented and tested.
- Public release approval is recorded before publishing.

Current result:

- GitHub Actions are pinned by commit SHA in CI and release workflows.
- GoReleaser is pinned to `v2.16.0` in the release workflow.
- GoReleaser signs the `SHA256SUMS` file with GPG and publishes the detached `.sig` asset.
- Terraform Registry manifest is added at `terraform-registry-manifest.json`, included in checksums, and published as a release asset.
- GitHub provenance attestations are generated for release artifacts.
- Secret scanning is added to CI and release workflows through `make secret-scan`.
- Public publication, GPG creation, Vault storage, GitHub environment secrets, Terraform Registry key setup, rotation, and revocation are documented in `docs/PUBLICATION.md`.
- Manual gates remain outside the repository: make GitHub repository public, create/store the GPG key in Vault, configure GitHub release environment secrets, add the public key to Terraform Registry, and publish the provider from the Terraform Registry UI.

## Risks

Open decisions are tracked in [DECISIONS.md](DECISIONS.md).

Known risks:

- Mailu API surface may not cover all admin UI features.
- API responses may not expose enough stable identifiers.
- Mailu may normalize email addresses differently than Terraform configuration.
- Password fields are inherently write-only, requiring explicit Terraform state handling.

Contribution and review workflow is documented in [CONTRIBUTING.md](CONTRIBUTING.md).
