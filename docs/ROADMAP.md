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
- Install Go 1.24+ locally.
- Confirm Terraform 1.8+ is available.
- Define local test variables with `MAILU_ENDPOINT`, `MAILU_API_TOKEN`, and `MAILU_ACC_DOMAIN`.
- Store local secrets outside Git.

Exit criteria:

- README links to project, Terraform, contribution, decision, API, and roadmap docs.
- Terraform provider documentation exists under `docs/index.md`.
- Planned resource and data source pages exist under `docs/resources` and `docs/data-sources`.
- `.env.local` and other local secrets are ignored by Git.
- `go version` reports Go 1.24+.
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
  name = "cetic-group.com"
}

resource "mailu_user" "admin" {
  email        = "admin@cetic-group.com"
  raw_password = var.admin_password
  quota_bytes  = 1073741824
  enabled      = true
  global_admin = true
}

resource "mailu_alias" "postmaster" {
  email       = "postmaster@cetic-group.com"
  destination = ["admin@cetic-group.com"]
}
```

Import IDs:

- Domain: `cetic-group.com`
- User: `admin@cetic-group.com`
- Alias: `postmaster@cetic-group.com`

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

Status: pending.

Goal: document how this provider composes with DNS providers.

Tasks:

- Produce examples for MX, SPF, DKIM, DMARC, MTA-STS, and autoconfig records.
- Add an IONOS-oriented example if CETIC Group continues using IONOS DNS automation.
- Document ownership boundaries: Mailu provider manages Mailu state, DNS provider manages DNS state.
- Add an example that consumes `mailu_dkim` or domain DNS output in DNS records, if implemented.
- Document which DNS records are mandatory, recommended, or optional.

Exit criteria:

- A new domain can be onboarded with Mailu and DNS using Terraform examples.
- DNS examples do not require secrets in source files.

## Phase 8 - Release Engineering

Status: pending.

Goal: make releases repeatable and installable.

Tasks:

- Add GitHub Actions CI.
- Add `goreleaser`.
- Generate provider docs with `terraform-plugin-docs`.
- Keep Terraform Registry documentation under `docs/index.md`, `docs/resources`, and `docs/data-sources`.
- Add changelog.
- Add semantic versioning.
- Decide public Terraform Registry vs internal mirror.
- Configure release signing or checksum publication.
- Add release notes template.
- Add upgrade notes for breaking changes.
- Publish `v0.1.0` as internal or public provider depending on CETIC Group policy.
- Run release readiness review.

Exit criteria:

- Tagged releases publish checksums and platform binaries.
- The provider can be installed from the Terraform Registry or an internal mirror.
- Documentation generated for the release matches implemented schemas.
- Release process is repeatable from a clean checkout.

## Phase 9 - Hardening

Status: pending.

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

## Phase 10 - Production Adoption

Status: pending.

Goal: migrate CETIC Group Mailu management to Terraform safely.

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

- Existing production objects are imported without unintended changes.
- First production apply succeeds.
- Operational runbook exists.

## Risks

Open decisions are tracked in [DECISIONS.md](DECISIONS.md).

Known risks:

- Mailu API surface may not cover all admin UI features.
- API responses may not expose enough stable identifiers.
- Mailu may normalize email addresses differently than Terraform configuration.
- Password fields are inherently write-only, requiring explicit Terraform state handling.

Contribution and review workflow is documented in [CONTRIBUTING.md](CONTRIBUTING.md).
