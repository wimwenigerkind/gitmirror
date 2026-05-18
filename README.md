# gitmirror

Mirrors repositories from git hosting providers to local bare repositories. First run does `git clone --mirror`; subsequent runs `git fetch --prune` to stay in sync.

Currently supports:
- **Bitbucket Cloud**

## Installation

### Docker

```bash
docker run --rm \
  -v $(pwd)/config.yaml:/work/config.yaml \
  -v $(pwd)/mirrors:/work/mirrors \
  ghcr.io/wimwenigerkind/gitmirror:latest
```

### Homebrew (macOS)

```bash
brew install --cask wimwenigerkind/tap/gitmirror
```

### Linux (systemd timer)

For a server install that mirrors on a schedule:

```bash
curl -fsSL https://raw.githubusercontent.com/wimwenigerkind/gitmirror/main/deploy/install.sh -o install.sh
less install.sh           # always inspect remote scripts before piping them to sudo
sudo bash install.sh
```

The script downloads the latest release, verifies the SHA256 against `checksums.txt`, installs the binary to `/usr/local/bin/gitmirror`, creates a `gitmirror` system user under `/var/lib/gitmirror`, and drops the service + timer units into `/etc/systemd/system/`. A placeholder `config.yaml` is created, edit it, then:

```bash
sudo systemctl enable --now gitmirror.timer
```

Pin a specific version: `VERSION=v1.2.3 sudo -E bash install.sh`.

Useful commands afterwards:

```bash
systemctl list-timers gitmirror.timer     # next scheduled run
systemctl start gitmirror.service         # ad-hoc run
journalctl -u gitmirror.service -f        # tail logs
```

### Binary

Download the archive for your OS/arch from the [releases page](https://github.com/wimwenigerkind/gitmirror/releases), extract, and put `gitmirror` on your `PATH`.

## Requirements

- `git` on PATH (not required when running via Docker)

## Configuration

Create `config.yaml`

```yaml
destination: ./mirrors    # optional, defaults to ./mirrors
concurrency: 4            # optional, defaults to 4

provider:
    provider-name:
        type: bitbucket
        owner: workspace # https://support.atlassian.com/bitbucket-cloud/docs/what-is-a-workspace/
        token: token # https://support.atlassian.com/bitbucket-cloud/docs/access-tokens/
```

Multiple providers can be listed under `provider:` and are processed sequentially; repos within a provider are mirrored in parallel up to `concurrency`.

## Not mirrored

- Git LFS objects (only refs and git objects are fetched)
- PRs, issues, wikis, server-side hooks, repository settings, permissions

## License

See [LICENSE](LICENSE).