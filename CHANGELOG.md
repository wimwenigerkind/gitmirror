# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.3] - 2026-05-18

### Added

- Project-aware on-disk layout. Bitbucket repositories are now mirrored to
  `{destination}/{provider}/{project-key}/{slug}.git`, grouped by their
  Bitbucket project. Providers without a project concept (e.g. future GitHub
  support) remain flat at `{destination}/{provider}/{slug}.git`.
- `Repository.Project` field on the provider interface to carry the optional
  sub-namespace from the source provider.

### Changed

- `BitbucketProvider.ListRepositories` now decodes `project.key` from the
  Bitbucket API response and propagates it through `Repository.Project`.

## [0.0.2] - 2026-05-18

### Fixed

- `deploy/install.sh`: version detection no longer aborts with `curl: (23)`
  when `grep -m1` closes the pipe early under `set -o pipefail`. The API
  response is now buffered into a variable before parsing.
- `deploy/install.sh`: checksum filename now matches the GoReleaser v2 default
  (`gitmirror_{version}_checksums.txt`) instead of the previously hardcoded
  `checksums.txt`.

## [0.0.1] - 2026-05-18

Initial release.

### Added

- Bitbucket Cloud provider with workspace-scoped repository listing, pagination
  via the `next` link, and Bearer-token authentication.
- `AuthenticatedURL` for Bitbucket, injecting `x-token-auth:{token}` into the
  HTTPS clone URL via `net/url`.
- Bare mirror sync (`internal/mirror`): `git clone --mirror` on first run,
  `git -C {dest} fetch --prune origin` on subsequent runs, with the remote URL
  refreshed each run to handle token rotation.
- YAML config loader (`internal/config`) with optional top-level `destination`
  (default `./mirrors`) and `concurrency` (default `4`).
- Provider factory dispatching on `type` (`bitbucket`, `github` stub).
- Per-provider parallel mirroring via `golang.org/x/sync/errgroup.SetLimit`.
- GoReleaser pipeline (`.goreleaser.yaml`) producing tar.gz/zip archives for
  Linux/macOS/Windows on amd64/arm64/i386.
- Docker image published to `ghcr.io/wimwenigerkind/gitmirror` via
  `dockers_v2`, multi-arch (`linux/amd64`, `linux/arm64`), based on
  `alpine:3.23` with `git` and `ca-certificates`.
- Homebrew cask published to `wimwenigerkind/homebrew-tap` with a post-install
  hook that clears the macOS quarantine attribute.
- GitHub Actions release workflow (`.github/workflows/release.yml`) wiring
  QEMU, Docker Buildx, GHCR login, and GoReleaser.
- systemd service + timer units in `deploy/systemd/` for scheduled mirroring
  with sandboxing (`ProtectSystem=strict`, `NoNewPrivileges=true`, dedicated
  `gitmirror` system user, `Persistent=true` to catch up after downtime).
- Linux installer script (`deploy/install.sh`) that downloads the latest
  release archive, verifies the SHA256 against the GoReleaser checksums file,
  installs the binary to `/usr/local/bin`, creates the `gitmirror` system
  user under `/var/lib/gitmirror`, and drops the systemd units into
  `/etc/systemd/system/`.

[unreleased]: https://github.com/wimwenigerkind/gitmirror/compare/v0.0.3...HEAD
[0.0.3]: https://github.com/wimwenigerkind/gitmirror/compare/v0.0.2...v0.0.3
[0.0.2]: https://github.com/wimwenigerkind/gitmirror/compare/v0.0.1...v0.0.2
[0.0.1]: https://github.com/wimwenigerkind/gitmirror/releases/tag/v0.0.1