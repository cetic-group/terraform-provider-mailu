# Generic DNS Composition Example

This example creates a Mailu domain and exposes DNS values that can be passed to the DNS provider that owns the zone.

It intentionally outputs DNS values instead of creating DNS records because DNS provider schemas differ.

Run:

```shell
terraform init
terraform plan \
  -var domain=example.com \
  -var tls_report_email=postmaster@example.com \
  -var mta_sts_policy_id=2026062901
```

The Mailu API token can be passed with `MAILU_API_TOKEN` or `var.mailu_api_token`.
