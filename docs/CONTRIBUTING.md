# Contributing

This project keeps documentation concise and implementation decisions explicit.

## Workflow

1. Confirm the Mailu API behavior before implementing a resource.
2. Update the resource model and decisions when behavior is clarified.
3. Implement the smallest coherent change.
4. Add unit tests, provider tests, and import coverage where applicable.
5. Run the required agent reviews before considering the change complete.

## Project Agents

The review agents live in [../../.agents](../../.agents).

Use them in this order:

1. Senior Developer Architect: architecture, resource model, Terraform Plugin Framework usage.
2. Senior QA: tests, acceptance criteria, regression risk.
3. Senior Application Security: secrets, authentication, Terraform state, destructive behavior.

All three are required for:

- New resources or data sources.
- Provider configuration changes.
- Authentication or token handling.
- Sensitive attributes.
- Import and state migration behavior.
- Destructive operations.
- Acceptance tests, CI, and releases.

## Quality Gates

Before release:

- `go test ./...`
- `terraform fmt -recursive`
- Acceptance tests in a controlled Mailu environment.
- Documentation examples match implemented schemas.
- Sensitive values are not present in examples, logs, test fixtures, or docs.

## Terraform Documentation

Provider documentation must remain compatible with Terraform Registry and `terraform-plugin-docs` conventions:

- Provider page: `docs/index.md`.
- Resource pages: `docs/resources/<name>.md`.
- Data source pages: `docs/data-sources/<name>.md`.
- Resource examples: `examples/resources/<type>/resource.tf`.
- Data source examples: `examples/data-sources/<type>/data-source.tf`.

Do not mark a resource as implemented in the Terraform documentation until the provider code, tests, import behavior, and examples are complete.
