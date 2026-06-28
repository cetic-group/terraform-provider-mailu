# mailu_server_info

Reads Mailu server metadata.

Status: blocked. No server metadata endpoint is exposed in the Swagger document.

## Example Usage

```terraform
data "mailu_server_info" "current" {}
```

## Schema

### Read-Only

- `version` (String) Mailu version, if exposed by the API.
- `hostname` (String) Mailu hostname, if exposed by the API.
- `features` (Set of String) Enabled Mailu features, if exposed by the API.
