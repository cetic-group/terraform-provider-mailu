# Mailu Provider

The Mailu provider manages Mailu domains, users, aliases, and related mail objects through the Mailu admin API.

Status: scaffolded. Resources and data sources are planned and will be enabled after the Mailu API contract is confirmed.

## Example Usage

```terraform
terraform {
  required_providers {
    mailu = {
      source  = "cetic-group/mailu"
      version = "0.1.0"
    }
  }
}

provider "mailu" {
  endpoint = "https://mail.cetic-group.com/api/v1"
  token    = var.mailu_api_token
}
```

## Schema

### Optional

- `endpoint` (String) Mailu API endpoint, for example `https://mail.example.com/api/v1`. Can also be set with `MAILU_ENDPOINT`.
- `token` (String, Sensitive) Mailu API token. Can also be set with `MAILU_API_TOKEN`.

## Environment Variables

- `MAILU_ENDPOINT`
- `MAILU_API_TOKEN`

## Planned Resources

- `mailu_domain`
- `mailu_user`
- `mailu_alias`
- `mailu_forward`
- `mailu_relay`
- `mailu_alternative_domain`
- `mailu_domain_manager`
- `mailu_token`

## Planned Data Sources

- `mailu_domain`
- `mailu_user`
- `mailu_dkim`

## Deferred

- `mailu_fetchmail`: no Swagger endpoint found.
- `mailu_server_info`: no Swagger endpoint found.
