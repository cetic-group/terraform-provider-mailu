# Import Examples

Import commands:

```shell
terraform import mailu_domain.example example.com
terraform import mailu_user.admin admin@example.com
terraform import mailu_alias.postmaster postmaster@example.com
terraform import mailu_alternative_domain.legacy legacy.example.com
terraform import mailu_domain_manager.admin example.com/admin@example.com
terraform import mailu_relay.example example.com
terraform import mailu_token.admin 42
```

Terraform import blocks are shown in `examples/production/imports.tf.example`.

For production adoption, follow [Production Adoption Runbook](../../docs/PRODUCTION_ADOPTION.md).
