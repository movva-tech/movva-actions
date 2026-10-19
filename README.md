# Movva Actions

Reusable GitHub Actions for a Python release workflow. The top-level composite action is intended to prepare a Python project for release: install dependencies, package the source, create a GitHub Release, and notify an external webhook. The repository also exposes two smaller actions for dependency scanning and release notifications.

## Actions

| Action | Purpose |
| --- | --- |
| `./` | Orchestrates the Python build, packaging, GitHub Release, and notification workflow. |
| `./blacklist-scanner` | Starts a Go-based scanner for a Python requirements file. |
| `./webhook-notifier` | Starts a Go-based notifier with release metadata and a bearer-token secret. |

## Release workflow

The root action (`action.yml`) performs these steps:

1. Sets up the requested Python version with pip caching.
2. Invokes the dependency scanner for `requirements.txt`.
3. Installs dependencies with `pip install -r requirements.txt`.
4. Creates `dist/project-bundle.zip` and `dist/project-bundle.tar.gz` from `src/`, `requirements.txt`, and `main.py`.
5. Creates a GitHub Release for the supplied tag.
6. Invokes the webhook notifier with the repository name and release tag.

It expects a conventional Python project layout containing `requirements.txt`, `main.py`, and a `src/` directory.

See the complete [release workflow example](example/release-example.yml).

### Intended workflow usage

Replace `OWNER/REPOSITORY@VERSION` with the repository and immutable version you publish.

```yaml
name: Release

on:
  push:
    tags: ['v*']

permissions:
  contents: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: OWNER/REPOSITORY@VERSION
        with:
          python-version: '3.12'
          tag-name: ${{ github.ref_name }}
          github-token: ${{ secrets.GITHUB_TOKEN }}
          webhook-secret: ${{ secrets.RELEASE_WEBHOOK_SECRET }}
```

## Dependency scanner

`blacklist-scanner` accepts one input:

| Input | Required | Default | Description |
| --- | --- | --- | --- |
| `requirements-file` | No | `requirements.txt` | Path to the Python requirements file to scan. |

The action installs Go 1.22 and exposes the selected path to the scanner as `INPUT_REQUIREMENTS_FILE`.

```yaml
- uses: OWNER/REPOSITORY/blacklist-scanner@VERSION
  with:
    requirements-file: requirements/production.txt
```

## Webhook notifier

`webhook-notifier` installs Go 1.22 and passes the following inputs to its notifier:

| Input | Required | Description |
| --- | --- | --- |
| `webhook-secret` | Yes | Bearer token used to authorize the request. Store it in GitHub Secrets. |
| `repository` | Yes | Repository name included in the notification. |
| `version` | Yes | Release version or tag included in the notification. |

```yaml
- uses: OWNER/REPOSITORY/webhook-notifier@VERSION
  with:
    webhook-secret: ${{ secrets.RELEASE_WEBHOOK_SECRET }}
    repository: ${{ github.repository }}
    version: ${{ github.ref_name }}
```

## Requirements and security

- Run the release action on a Linux runner with `zip` and `tar` available (for example, `ubuntu-latest`).
- Grant `contents: write` when the action creates a release with `GITHUB_TOKEN`.
- Pin third-party actions and this action to a trusted tag or commit SHA in production.
- Keep webhook secrets in GitHub Secrets; do not commit them to workflows.

## Current implementation status

This repository is an initial scaffold. The Go entry points are placeholders, so the scanner does not yet enforce a blacklist and the notifier does not yet send an HTTP request. The release configuration creates a release tag but does not currently attach the generated archives as release assets. Complete those implementations before relying on these actions in a production release process.
