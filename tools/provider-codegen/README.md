# provider-codegen

Code generation tooling for this provider.

## What it does

- Runs schema-driven artifact generation (`cmd/schema-to-artifacts`):
  - `data_source_model_gen.go`
  - `resource_model_gen.go`
  - `resource_diff_gen.go` (resources only)
- Runs post-generation steps (`go generate` in this module):
  - `copywrite license`
  - `copywrite headers`
  - `terraform fmt` for examples
  - `tfplugindocs generate`

## Usage

Run full pipeline (recommended), from repo root:

```bash
make generate
```

Run full pipeline directly:

```bash
cd tools/provider-codegen
go generate ./...
```

## Schema To Artifacts

`schema-to-artifacts` scans provider resource and data source schema files and writes generated files next to them.

Generated files:

- `data_source_model_gen.go`
- `resource_model_gen.go`
- `resource_diff_gen.go` (resources only)

Run schema artifact generation only:

```bash
cd tools/provider-codegen
go run ./cmd/schema-to-artifacts --root ../../internal/resources --schema-pattern resource_schema.go --schema-pattern data_source_schema.go
```

Useful flags:

- `--verbose` prints matched schema files and generated output paths.

## Codegen Annotations

You can control whether a schema attribute participates in generated
`HasRemoteChanges()` detection with a comment annotation above the map entry:

```go
// codegen:has_remote_changes=false
"desired_state": schema.StringAttribute{
    // ...
}
```

- Default behavior (no annotation): `codegen:has_remote_changes=true`
- Typical use case for `false`: provider-only control fields that should not trigger AWS update calls
