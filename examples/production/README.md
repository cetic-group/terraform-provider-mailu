# Production Adoption Example

This directory is a scaffold for production adoption.

It is intentionally incomplete and should not be applied as-is. Copy the `.example` files to real Terraform files only after the production backend, inventory, and review process are approved.

Recommended flow:

```shell
cp backend.tf.example backend.tf
cp imports.tf.example imports.tf
terraform init
terraform validate
terraform plan -refresh-only
terraform plan
```

Before production import, complete:

- `INVENTORY_TEMPLATE.md`
- `PLAN_REVIEW_CHECKLIST.md`

Use environment variables for Mailu credentials:

```shell
export MAILU_ENDPOINT="https://mail.example.com/api/v1"
export MAILU_API_TOKEN="..."
```

Stop immediately if `terraform plan` proposes a delete or an unexpected update.
