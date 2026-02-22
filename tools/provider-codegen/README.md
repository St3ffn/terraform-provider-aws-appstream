# provider-codegen

Tooling module for provider generation tasks.

## What it does

Runs the repository generation pipeline defined in `tools.go`:

1. Generate license files (`copywrite license`)
2. Generate/update file headers (`copywrite headers`)
3. Format Terraform examples (`terraform fmt`)
4. Generate provider docs (`tfplugindocs`)

## Usage

From repository root:

```bash
make generate
```

Or directly:

```bash
cd tools/provider-codegen
go generate ./...
```

## Requirements

- Go
- Terraform CLI (for `terraform fmt`)

## Notes

- Commands operate on repo root via relative paths from this module.
- `tools.go` is build-tagged with `generate` and intended for `go generate` only.
