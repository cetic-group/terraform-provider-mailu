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
- Alternative domains
- Domain managers
- Relays
- Authentication tokens
- Mailu API discovery and import workflows

## Requirements

- Go 1.25.8+
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
make fmt
make build
make test
make testacc
make install-local
```

The resources `mailu_domain`, `mailu_user`, `mailu_alias`, `mailu_alternative_domain`, `mailu_domain_manager`, `mailu_relay`, and `mailu_token` are implemented. The data sources `mailu_domain`, `mailu_user`, and `mailu_dkim` are implemented.

`make testacc` requires:

```shell
export TF_ACC=1
export MAILU_ENDPOINT="https://mail.cetic-group.com/api/v1"
export MAILU_API_TOKEN="..."
export MAILU_ACC_DOMAIN="example.com"
```

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
- [DNS integration](docs/DNS.md)
- [Release process](docs/RELEASE.md)
- [Private provider installation](docs/PRIVATE_INSTALL.md)
- [Hardening guide](docs/HARDENING.md)
- [Production adoption runbook](docs/PRODUCTION_ADOPTION.md)
- [Upgrade guide](docs/UPGRADING.md)
- [Changelog](CHANGELOG.md)
- [Contributing](docs/CONTRIBUTING.md)
- [Decisions](docs/DECISIONS.md)
