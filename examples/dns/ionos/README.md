# IONOS DNS-Oriented Example

This example reads Mailu DNS values and shapes them into a payload that can be mapped to the selected IONOS DNS automation layer.

No IONOS secret is stored in this repository. Keep the IONOS API key outside Git, following `mailu/mailu-data/ionos.env.example`.

Run:

```shell
terraform init
terraform plan \
  -var domain=example.com \
  -var tls_report_email=postmaster@example.com \
  -var mta_sts_policy_id=2026062901
```

Map `output.ionos_record_payload` to the IONOS DNS mechanism used by the environment.
