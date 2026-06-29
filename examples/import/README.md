# Import Examples

Import commands:

```shell
terraform import mailu_domain.cetic cetic-group.com
terraform import mailu_user.admin admin@cetic-group.com
terraform import mailu_alias.postmaster postmaster@cetic-group.com
terraform import mailu_alternative_domain.legacy legacy-cetic-group.com
terraform import mailu_domain_manager.admin cetic-group.com/admin@cetic-group.com
terraform import mailu_relay.cetic cetic-group.com
terraform import mailu_token.admin 42
```

Terraform import blocks are shown in `examples/production/imports.tf.example`.

For production adoption, follow [Production Adoption Runbook](../../docs/PRODUCTION_ADOPTION.md).
