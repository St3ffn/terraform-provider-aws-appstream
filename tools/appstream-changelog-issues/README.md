# appstream-changelog-issues

Checks AWS AppStream SDK changelog updates and creates/updates one tracking GitHub issue.

## What it does

1. Reads current `github.com/aws/aws-sdk-go-v2/service/appstream` version from provider `go.mod`.
2. Fetches and parses AppStream changelog from AWS SDK v2.
3. Keeps releases newer than current version.
4. Keeps only releases that contain at least one `**Feature**:` note.
5. Upserts a single issue in your repo (marker + title), with labels:
   - `dependencies`
   - `appstream-sdk-watch`

## Defaults

- Changelog URL: `https://raw.githubusercontent.com/aws/aws-sdk-go-v2/main/service/appstream/CHANGELOG.md`
- Issue title: `AWS SDK service/appstream feature updates available`
- Issue marker: `<!-- appstream-changelog-issues -->`

## Usage

From repo root:

```bash
cd tools/appstream-changelog-issues
go run ./cmd/appstream-changelog-issues
```

Dry run:

```bash
cd tools/appstream-changelog-issues
go run ./cmd/appstream-changelog-issues --dry-run
```

## Required environment

- `GITHUB_REPOSITORY` in `owner/repo` format (or pass `--repo`)
- `GITHUB_TOKEN` for write mode (`--dry-run=false`)

## CLI flags

- `--provider-go-mod` default: `go.mod`
- `--repo` default: from `GITHUB_REPOSITORY`
- `--github-token` default: from `GITHUB_TOKEN`
- `--dry-run` default: `false`
- `--timeout` default: `5m`
