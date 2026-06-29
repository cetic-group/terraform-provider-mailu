# Production Adoption Runbook

This runbook describes how operators can migrate existing Mailu objects to Terraform without unintended changes.

Do not run production `apply` from this repository until the inventory, imports, backend, and review gates below are complete.

Current phase status: Phase 10 is complete for minimal adoption validation. The tested Mailu installation only contained an initial domain and admin user, and both were imported and refreshed successfully with a no-op plan. A remote production backend, broader object rollout, and first low-risk production apply are deferred until an operator starts a wider production rollout.

## Scope

Phase 10 covers production adoption:

- Inventory existing Mailu objects.
- Write Terraform configuration that mirrors current Mailu state.
- Import existing objects into a non-production state first.
- Compare Terraform state with Mailu API reads.
- Prepare the production backend and access policy.
- Produce a no-op production plan before any change.
- Apply one low-risk production change after review.

Public Terraform Registry publication and release supply-chain hardening are handled in Phase 11.

Use provider release `0.1.0-rc.2` or newer for production adoption. Do not use `0.1.0-rc.1`, which was published before the phase 9 hardening changes.

## Required Inputs

- Mailu API endpoint.
- Admin API token with enough read permission for inventory.
- Approved Terraform backend location.
- List of production domains to manage.
- Reviewers for architecture, QA, and application security.

Store secrets only in environment variables or an approved secret manager:

```shell
export MAILU_ENDPOINT="https://mail.example.com/api/v1"
export MAILU_API_TOKEN="..."
```

Never commit `.tfvars`, state files, plans, or API output containing secrets.

## Step 1 - Inventory

Collect the current Mailu objects before writing Terraform resources:

- Domains.
- Users.
- Aliases.
- Alternative domains.
- Domain managers.
- Relays.
- Tokens, metadata only.
- DKIM/DNS values for comparison.

Recommended inventory format:

```text
domain: example.com
users:
  - admin@example.com
aliases:
  - postmaster@example.com
domain_managers:
  - example.com/admin@example.com
```

Keep the inventory in a private operational location. Do not commit production user lists unless the organization explicitly approves it.

Use `examples/production/INVENTORY_TEMPLATE.md` as the starting format.

## Step 2 - Write Matching Terraform Configuration

Create Terraform resources that match the current Mailu state. Do not introduce desired future changes in this step.

Rules:

- Use lowercase domain names and email addresses.
- Use Terraform references to express dependencies.
- Do not configure `mailu_relay.smtp` with embedded credentials.
- Do not expect `mailu_token.token` to be available after import or read.
- Keep DNS records in DNS provider configurations, not in this provider.

## Step 3 - Import Into Non-Production State First

Use a temporary validation workspace or local encrypted state to test imports.

Import IDs:

```shell
terraform import mailu_domain.example example.com
terraform import mailu_user.admin admin@example.com
terraform import mailu_alias.postmaster postmaster@example.com
terraform import mailu_alternative_domain.old old-example.com
terraform import mailu_domain_manager.admin example.com/admin@example.com
terraform import mailu_relay.example example.com
terraform import mailu_token.admin 42
```

After imports:

```shell
terraform state list
terraform plan -refresh-only
terraform plan
```

The expected result before production adoption is a no-op `terraform plan`.

Use `examples/production/PLAN_REVIEW_CHECKLIST.md` to record the review result.

## Minimal Adoption Validation

On 2026-06-29, a minimal non-production adoption validation was completed from:

```text
/private/tmp/mailu-terraform-adoption
```

Imported objects:

- `mailu_domain.example` with ID `example.com`.
- `mailu_user.admin` with ID `admin@example.com`.

Validation command:

```shell
terraform plan -refresh=true -detailed-exitcode
```

Observed result:

```text
mailu_domain.example: Refreshing state... [id=example.com]
mailu_user.admin: Refreshing state... [id=admin@example.com]

No changes. Your infrastructure matches the configuration.
```

Conclusion:

- The provider can refresh the imported Mailu domain and admin user.
- The Terraform configuration matches the current minimal Mailu state.
- No delete, replacement, or unexpected update was proposed.
- This validates the minimal non-production import path.

Deferred before wider production rollout:

- Approved production backend configuration.
- Production import into the approved backend.
- Production no-op plan review.
- First low-risk production apply.
- Post-apply no-op plan.

## Step 4 - Production Backend

Before production import, configure a remote backend with:

- Encryption at rest.
- Strict write access.
- State locking.
- Audit logs.
- Backup and restore procedure.
- Access limited to the Mailu operations group.

Do not use local production state for the final migration.

## Step 5 - Production Import

Run production imports only after non-production import validation is complete.

Recommended sequence:

1. Domains.
2. Users.
3. Aliases.
4. Alternative domains.
5. Domain managers.
6. Relays.
7. Tokens metadata.
8. Data source verification.

After every batch:

```shell
terraform plan -refresh-only
terraform plan
```

Stop if Terraform proposes a delete or unexpected update.

## Step 6 - First Production Change

The first production change must be low risk. Recommended candidates:

- Update a non-critical comment.
- Add a temporary test alias under an approved domain.
- Adjust metadata on a disposable relay entry.

Do not use the first apply to change passwords, quotas, managers, relays, or aliases used by production mail flow.

## Rollback

Rollback depends on the operation:

- If an import is wrong, remove the Terraform state address with `terraform state rm`; this does not delete Mailu objects.
- If a plan proposes unintended changes, do not apply; fix configuration or state first.
- If an apply changes metadata incorrectly, restore the previous Terraform configuration and apply after review.
- If an apply impacts mail routing, pause Terraform usage and restore Mailu settings manually through the Mailu admin UI/API.

Always record the exact command, operator, timestamp, state version, and plan reviewed.

## Approval Gates

Phase 10 cannot be considered fully executed until:

- Inventory is complete.
- Non-production import is validated.
- Production backend is approved.
- Production import produces a no-op plan.
- Architecture, QA, and application security reviews approve the plan.
- First low-risk production apply succeeds.
- The operational runbook is accepted by the responsible operations team.
