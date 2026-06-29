# Mailu Provider

The Mailu provider manages Mailu domains, users, aliases, and related mail objects through the Mailu admin API.

Status: MVP resources and data sources are implemented for domains, users, and aliases.

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
- `timeout_seconds` (Number) HTTP client timeout in seconds. Can also be set with `MAILU_TIMEOUT_SECONDS`.
- `max_retries` (Number) Maximum retry count for retryable API responses. Can also be set with `MAILU_MAX_RETRIES`.
- `user_agent` (String) User-Agent sent to the Mailu API. Can also be set with `MAILU_USER_AGENT`.
- `insecure_skip_tls_verify` (Boolean) Skip TLS certificate verification. Can also be set with `MAILU_INSECURE_SKIP_TLS_VERIFY`; intended only for lab environments.

## Environment Variables

- `MAILU_ENDPOINT`
- `MAILU_API_TOKEN`
- `MAILU_TIMEOUT_SECONDS`
- `MAILU_MAX_RETRIES`
- `MAILU_USER_AGENT`
- `MAILU_INSECURE_SKIP_TLS_VERIFY`

## Resources

- `mailu_domain`
- `mailu_user`
- `mailu_alias`

## Planned Resources

- `mailu_forward`
- `mailu_relay`
- `mailu_alternative_domain`
- `mailu_domain_manager`
- `mailu_token`

## Data Sources

- `mailu_domain`
- `mailu_user`

## Planned Data Sources

- `mailu_dkim`

## Deferred

- `mailu_fetchmail`: no Swagger endpoint found.
- `mailu_server_info`: no Swagger endpoint found.
