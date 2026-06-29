# Production Plan Review Checklist

Use this checklist before every production import batch and before the first production apply.

## Required Commands

```shell
terraform init
terraform validate
terraform plan -refresh-only -detailed-exitcode
terraform plan -detailed-exitcode
```

Expected result after import:

- `terraform plan -refresh-only -detailed-exitcode` returns `0`.
- `terraform plan -detailed-exitcode` returns `0`.

Exit code `2` is allowed only when every proposed change is classified and approved. Exit code `1` blocks the rollout.

## Blocking Changes

Stop the rollout if the plan contains:

- Any delete.
- Any replacement of a domain, user, alias, alternative domain, manager, relay, or token.
- Any unexpected change to `global_admin`, managers, relays, aliases, forwarding, or quota.
- Any password rotation not explicitly approved.
- Any diff caused by incomplete inventory.
- Any object absent from the production inventory.

## Review Record

- Plan command:
- Plan timestamp:
- Terraform version:
- Provider version:
- Backend:
- State version:
- Reviewer architecture:
- Reviewer QA:
- Reviewer security:
- Approved for import:
- Approved for first apply:

## First Apply Candidate

The first production apply must be low risk.

Approved candidates:

- Metadata comment update on a non-critical object.
- Temporary alias under an approved low-risk domain.
- Configuration-only correction with no mail routing impact.

Rejected candidates:

- User deletion.
- Domain deletion.
- Password rotation.
- Mail routing change.
- Relay credential change.
- Manager/global admin change.
