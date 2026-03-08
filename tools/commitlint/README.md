# commitlint

Go-native commitlint tool for this repository.

It validates commit messages against the [Conventional Commits 1.0.0 specification](https://www.conventionalcommits.org/en/v1.0.0/#specification) using `go-git` for commit selection.

## Usage

From repository root:

```bash
cd tools/commitlint
go run ./cmd/commitlint
```

## Commit selection

Selection flags are mutually exclusive.

- default (no selector): current branch diff against `--base` (`main` by default)
  - equivalent to: `--branch-diff <current-branch> --base <base>`
- `--all`: all commits in repository history
- `--branch <name>`: all commits reachable from branch tip
- `--branch-diff <name>`: commits in `base..branch`
- `--from <sha> --to <sha>`: inclusive commit range

## Output formats

- `--format text` (default)
- `--format json`
- `--format github` (GitHub workflow annotations)
- `--format markdown` (useful for PR comments)
- `--format rdjsonl` (reviewdog-compatible diagnostics)
- `--fail-level error|none` (default: `error`)
  - `error`: exit non-zero when violations are found
  - `none`: always exit zero for lint findings (tool/runtime errors still exit non-zero)

Color controls for text output:

- `--color auto|always|never` (default: `auto`)
- `--no-color`

## Environment variables

All flags can be configured via environment variables with the `COMMITLINT_` prefix.

Examples:

- `COMMITLINT_BASE=origin/main`
- `COMMITLINT_FORMAT=github`
- `COMMITLINT_NO_COLOR=true`
- `COMMITLINT_INCLUDE_MERGE_COMMITS=true`

## Useful examples

Default PR-style selection:

```bash
go run ./cmd/commitlint
```

Check branch diff explicitly:

```bash
go run ./cmd/commitlint --branch-diff feat/my-branch --base main
```

Check range:

```bash
go run ./cmd/commitlint --from <base_sha> --to <head_sha>
```

GitHub annotation output:

```bash
go run ./cmd/commitlint --format github --no-color
```

reviewdog output:

```bash
go run ./cmd/commitlint --format rdjsonl --no-color
```

Generate completion:

```bash
go run ./cmd/commitlint completion zsh > "${fpath[1]}/_commitlint"
```
