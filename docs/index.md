# Mailu Provider

The Mailu provider manages Mailu domains, users, aliases, and related mail objects through the Mailu admin API.

Status: MVP and extended mail resources are implemented for domains, users, aliases, alternative domains, domain managers, relays, tokens, and DKIM metadata.

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
- `mailu_relay`
- `mailu_alternative_domain`
- `mailu_domain_manager`
- `mailu_token`

## Modeled On Existing Resources

- User forwarding is managed through `mailu_user` fields: `forward_enabled`, `forward_destination`, and `forward_keep`.

## Data Sources

- `mailu_domain`
- `mailu_user`
- `mailu_dkim`

## DNS Integration

DNS records are managed by DNS providers, not by this provider. See [DNS integration](DNS.md) for MX, SPF, DKIM, DMARC, MTA-STS, autoconfig, TLSA, and IONOS-oriented patterns.

## Deferred

- `mailu_fetchmail`: no Swagger endpoint found.
- `mailu_server_info`: no Swagger endpoint found.
