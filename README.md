# Terraform Provider Mailu

Terraform provider for managing Mailu resources for CETIC Group.

Provider address:

```hcl
cetic-group/mailu
```

Initial scope:

- Domains
- Users
- Aliases
- Forwards
- DKIM metadata
- Mailu API discovery and import workflows

## Requirements

- Go 1.24+
- Terraform 1.8+
- A Mailu instance with the admin API enabled

The current CETIC Group Mailu installation already has these relevant settings in `mailu/.env`:

```env
API=true
WEB_API=/api
```

## Provider Configuration

```hcl
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

Environment variables are also supported:

```shell
export MAILU_ENDPOINT="https://mail.cetic-group.com/api/v1"
export MAILU_API_TOKEN="..."
```

## Development

```shell
go mod tidy
make test
make build
make install-local
```

Resources are intentionally not implemented until the Mailu API contract is captured from the running instance.

## Local Terraform Test

```shell
cd examples/basic
terraform init
terraform plan
```

## Documentation

- [Terraform provider documentation](docs/index.md)
- [Roadmap](docs/ROADMAP.md)
- [API contract](docs/API.md)
- [Resource model](docs/RESOURCE_MODEL.md)
- [Contributing](docs/CONTRIBUTING.md)
- [Decisions](docs/DECISIONS.md)
