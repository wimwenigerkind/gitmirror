#!/usr/bin/env bash
# Install gitmirror on a Linux host with a systemd timer.
# Usage: sudo bash install.sh   (optionally: VERSION=v1.2.3 sudo -E bash install.sh)
set -euo pipefail

REPO="wimwenigerkind/gitmirror"
PREFIX=/usr/local/bin
DATA_DIR=/var/lib/gitmirror
UNIT_DIR=/etc/systemd/system
USER_NAME=gitmirror

log() { printf '==> %s\n' "$*"; }
die() { printf 'error: %s\n' "$*" >&2; exit 1; }

[[ $EUID -eq 0 ]] || die "must run as root (sudo bash $0)"

for cmd in curl tar sha256sum systemctl install useradd; do
  command -v "$cmd" >/dev/null || die "missing required command: $cmd"
done
command -v git >/dev/null || die "git is required at runtime — install it (e.g. apt-get install git)"

case "$(uname -m)" in
  x86_64|amd64) ARCH=x86_64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  *) die "unsupported architecture: $(uname -m)" ;;
esac

VERSION="${VERSION:-}"
if [[ -z "$VERSION" ]]; then
  VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep -m1 '"tag_name":' \
    | sed -E 's/.*"tag_name":[[:space:]]*"([^"]+)".*/\1/')
fi
[[ -n "$VERSION" ]] || die "could not determine latest version"

log "Installing gitmirror ${VERSION} (Linux/${ARCH})"

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

ARCHIVE="gitmirror_Linux_${ARCH}.tar.gz"
BASE="https://github.com/${REPO}/releases/download/${VERSION}"

log "Downloading archive + checksums"
curl -fsSL --proto '=https' --tlsv1.2 "${BASE}/${ARCHIVE}"    -o "${TMP}/${ARCHIVE}"
curl -fsSL --proto '=https' --tlsv1.2 "${BASE}/checksums.txt" -o "${TMP}/checksums.txt"

log "Verifying SHA256"
(cd "$TMP" && sha256sum --ignore-missing -c checksums.txt >/dev/null) \
  || die "checksum verification failed"

log "Extracting"
tar -xzf "${TMP}/${ARCHIVE}" -C "$TMP" gitmirror

log "Installing binary to ${PREFIX}/gitmirror"
install -m 0755 -o root -g root "${TMP}/gitmirror" "${PREFIX}/gitmirror"

if ! id "$USER_NAME" &>/dev/null; then
  log "Creating system user ${USER_NAME}"
  useradd --system --create-home --home-dir "$DATA_DIR" \
          --shell /usr/sbin/nologin "$USER_NAME"
fi

log "Installing systemd units"
for unit in gitmirror.service gitmirror.timer; do
  curl -fsSL --proto '=https' --tlsv1.2 \
    "https://raw.githubusercontent.com/${REPO}/${VERSION}/deploy/systemd/${unit}" \
    -o "${UNIT_DIR}/${unit}"
  chmod 0644 "${UNIT_DIR}/${unit}"
  chown root:root "${UNIT_DIR}/${unit}"
done

if [[ ! -e "${DATA_DIR}/config.yaml" ]]; then
  log "Creating placeholder ${DATA_DIR}/config.yaml"
  cat > "${DATA_DIR}/config.yaml" <<'EOF'
destination: ./mirrors
concurrency: 4

provider:
EOF
  chown "${USER_NAME}:${USER_NAME}" "${DATA_DIR}/config.yaml"
  chmod 0640 "${DATA_DIR}/config.yaml"
fi

systemctl daemon-reload

cat <<EOF

Installed. Next steps:
  1) Edit ${DATA_DIR}/config.yaml and add your provider(s)
  2) Enable the timer:         systemctl enable --now gitmirror.timer
  3) Check next scheduled run: systemctl list-timers gitmirror.timer
  4) Tail logs:                journalctl -u gitmirror.service -f
  5) Ad-hoc run:               systemctl start gitmirror.service
EOF