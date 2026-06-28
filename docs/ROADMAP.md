# Roadmap

This roadmap tracks implementation of the `cetic-group/mailu` Terraform provider.

Process rule: phases must be completed in order. A phase is complete only when its exit criteria are satisfied and the required architecture, QA, and security reviews have no release-blocking findings.

## Phase -1 - Prerequisites

Status: pending.

Goal: make the development environment and target Mailu environment ready.

Tasks:

- Install Go 1.22+.
- Confirm Terraform 1.8+ is available.
- Run `go mod tidy`.
- Run `gofmt -w .`.
- Run `go test ./...`.
- Define the Mailu test endpoint.
- Create or identify a disposable test domain.
- Create a Mailu API token dedicated to provider development.
- Store local secrets outside Git using environment variables.
- Confirm no production mailbox/domain will be used by acceptance tests.

Exit criteria:

- `go version` reports Go 1.22+.
- `terraform version` reports Terraform 1.8+.
- `go test ./...` passes on the scaffold.
- `MAILU_ENDPOINT`, `MAILU_API_TOKEN`, and `MAILU_ACC_DOMAIN` are defined for local acceptance testing.
- Test environment ownership is documented.

## Phase 0 - Project Baseline

Status: started.

Deliverables:

- Provider repository scaffold.
- Go module `github.com/cetic-group/terraform-provider-mailu`.
- Provider address `registry.terraform.io/cetic-group/mailu`.
- Provider configuration through HCL and environment variables.
- Shared Mailu HTTP client.
- Initial examples, documentation, and project agents.

Exit criteria:

- `go test ./...` passes.
- Provider can be installed locally.
- Terraform can load the provider with `source = "cetic-group/mailu"`.
- README links to project, Terraform, contribution, decision, API, and roadmap docs.
- Terraform provider documentation exists under `docs/index.md`.
- Planned resource and data source pages exist under `docs/resources` and `docs/data-sources`.

## Phase 1 - Mailu API Discovery

Status: partially complete; blocked on runtime validation with a development token and disposable domain.

Goal: capture the real API contract from the CETIC Group Mailu instance before implementing Terraform resources.

Tasks:

- Confirm the exact authentication mechanism.
- Confirm endpoint paths for domains, users, aliases, DKIM, forwards, fetchmail, and server metadata.
- Capture list, read, create, update, and delete payloads.
- Capture response payloads for successful and failed requests.
- Capture pagination behavior, if any.
- Capture stable identifiers and natural import IDs.
- Confirm whether each object can be read after creation.
- Confirm whether updates are patch, replace, or delete/recreate operations.
- Confirm delete behavior for every object.
- Confirm password handling and whether password values are ever returned.
- Confirm whether admin flags, quota, enabled state, display name, spam threshold, and destinations are exposed.
- Confirm Mailu version and API compatibility assumptions.
- Write redacted examples in `docs/API.md`.
- Update `docs/RESOURCE_MODEL.md` with confirmed fields and unsupported fields.
- Update `docs/DECISIONS.md` with resolved authentication, delete, DNS, and password decisions.
- Run architecture, QA, and security reviews.

Exit criteria:

- Every planned Terraform resource has an endpoint mapping.
- Import IDs are specified.
- API gaps are documented with a fallback decision.
- Unsupported resources are explicitly marked as blocked or deferred.
- No implementation proceeds from unverified endpoint assumptions.

Current blocker:

- Swagger discovery is complete, but authenticated list/read/create/update/delete behavior still needs validation with `MAILU_API_TOKEN` and `MAILU_ACC_DOMAIN`.

## Phase 1.5 - Design Freeze For MVP

Status: pending.

Goal: freeze the first implementable provider surface before coding resources.

Tasks:

- Select MVP resources from confirmed API capabilities.
- Finalize schemas for `mailu_domain`, `mailu_user`, and `mailu_alias`.
- Mark each schema attribute as required, optional, computed, sensitive, or replacement-forcing.
- Define Terraform IDs and import ID formats.
- Define normalization rules for domains, localparts, and email addresses.
- Define password update behavior.
- Define delete strategy.
- Define drift behavior for every MVP resource.
- Define acceptance test fixtures and cleanup rules.
- Update Terraform documentation pages from planned schemas.

Exit criteria:

- `docs/RESOURCE_MODEL.md` matches confirmed API behavior.
- `docs/resources/*.md` matches the MVP schemas.
- `docs/data-sources/*.md` matches the MVP schemas.
- `docs/DECISIONS.md` has no open blocker for MVP implementation.
- Architecture, QA, and security reviews approve the MVP design.

## Phase 2 - Provider Foundation

Status: pending.

Tasks:

- Add client tests with `httptest`.
- Implement API error types.
- Add retry and timeout configuration.
- Add Terraform diagnostics helpers.
- Add provider-level `insecure_skip_tls_verify` only if needed for lab environments.
- Add user agent with provider version.
- Add acceptance test harness gated behind environment variables.
- Add redaction helpers for diagnostics and logs.
- Add environment-variable based test configuration.
- Add Makefile targets for `fmt`, `test`, and acceptance tests.
- Add local provider install instructions for the current platform.

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

## Phase 3 - MVP Resources

Status: pending.

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

## Phase 4 - Mail Operations

Status: pending.

Resources:

- `mailu_forward`
- `mailu_relay`
- `mailu_fetchmail`, only if a future API exposes it
- `mailu_token`, if Mailu exposes manageable application tokens.

Data sources:

- `mailu_dkim`
- `mailu_server_info`, only if a future API exposes it

Exit criteria:

- Common mailbox lifecycle is fully managed as code.
- DNS modules can consume DKIM values from data sources.
- Each added object follows the same discovery, design, test, documentation, and review process as MVP resources.

## Phase 5 - DNS Integration Patterns

Status: pending.

Goal: document how this provider composes with DNS providers.

Tasks:

- Produce examples for MX, SPF, DKIM, DMARC, MTA-STS, and autoconfig records.
- Add an IONOS-oriented example if CETIC Group continues using IONOS DNS automation.
- Document ownership boundaries: Mailu provider manages Mailu state, DNS provider manages DNS state.
- Add an example that consumes `mailu_dkim` output in DNS records, if the data source is implemented.
- Document which DNS records are mandatory, recommended, or optional.

Exit criteria:

- A new domain can be onboarded with Mailu and DNS using Terraform examples.
- DNS examples do not require secrets in source files.

## Phase 6 - Release Engineering

Status: pending.

Tasks:

- Add GitHub Actions CI.
- Add `goreleaser`.
- Generate provider docs with `terraform-plugin-docs`.
- Keep Terraform Registry documentation under `docs/index.md`, `docs/resources`, and `docs/data-sources`.
- Add changelog.
- Add semantic versioning.
- Publish `v0.1.0` as internal or public provider depending on CETIC Group policy.
- Run release readiness review.
- Decide public Terraform Registry vs internal mirror.
- Configure release signing or checksum publication.
- Add release notes template.
- Add upgrade notes for breaking changes.

Exit criteria:

- Tagged releases publish checksums and platform binaries.
- The provider can be installed from the Terraform Registry or an internal mirror.
- Documentation generated for the release matches implemented schemas.
- Release process is repeatable from a clean checkout.

## Phase 7 - Hardening

Status: pending.

Tasks:

- Rate limit handling.
- Better diff suppression for normalized email addresses.
- Case normalization policy for domains and localparts.
- State migration tests.
- Import validation.
- Partial-state handling for create/update failures.
- Detailed acceptance tests against disposable domains.
- Concurrency behavior review.
- Timeout and retry tuning against real Mailu behavior.
- Negative tests for authorization failures.
- State upgrade framework for future schema versions.
- Security review of CI logs and release artifacts.

Exit criteria:

- Provider behavior is stable enough for production mailbox/domain management.
- Production rollout checklist exists.
- Known limitations are documented in Terraform provider docs.

## Phase 8 - Production Adoption

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
