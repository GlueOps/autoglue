#!/usr/bin/env bash
set -euo pipefail

USER_TO_ADD=""
NONINTERACTIVE=0

print_usage() {
  cat <<'USAGE'
bootstrap_host.sh - install Docker and target prerequisites for Ansible.

Options:
  -u, --user USER   Add USER to the "docker" group (enables docker without sudo)
  -y, --yes         Run noninteractively, auto-accept prompts
  -h, --help        Show this help

Examples:
  sudo ./bootstrap_host.sh -u deploy
  sudo ./bootstrap_host.sh -y
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    -u|--user)
      USER_TO_ADD="${2:-}"
      shift 2
      ;;
    -y|--yes)
      NONINTERACTIVE=1
      shift
      ;;
    -h|--help)
      print_usage; exit 0
      ;;
    *)
      echo "Unknown option: $1" >&2
      print_usage
      exit 1
      ;;
  esac
done

if [[ "$(id -u)" -ne 0 ]]; then
  echo "Please run as root (use sudo)." >&2
  exit 1
fi

ID=""
ID_LIKE=""
if [[ -r /etc/os-release ]]; then
  # shellcheck disable=SC1091
  . /etc/os-release
fi

ID="${ID:-}"
ID_LIKE="${ID_LIKE:-}"

echo "Detected OS: ID='${ID}' ID_LIKE='${ID_LIKE}'"

apt_install() {
  apt-get update -y
  apt-get install -y --no-install-recommends "$@"
}

yum_or_dnf_install() {
  if command -v dnf >/dev/null 2>&1; then
    dnf install -y "$@"
  else
    yum install -y "$@"
  fi
}

apk_install() {
  apk add --no-cache "$@"
}

zypper_install() {
  zypper --non-interactive install -y "$@"
}

# Ensure basic prerequisites (sudo, curl, python3)
install_base_prereqs() {
  if [[ "$ID" == "alpine" ]]; then
    apk update
    apk_install sudo curl python3 py3-pip shadow
  elif [[ "$ID" == "debian" || "$ID" == "ubuntu" || "$ID_LIKE" == *"debian"* ]]; then
    apt_install sudo curl python3 python3-venv python3-pip ca-certificates gnupg lsb-release
  elif [[ "$ID" == "amzn" || "$ID_LIKE" == *"rhel"* || "$ID" == "fedora" || "$ID_LIKE" == *"fedora"* ]]; then
    yum_or_dnf_install sudo curl python3 python3-pip ca-certificates
  elif [[ "$ID" == "opensuse-leap" || "$ID_LIKE" == *"suse"* ]]; then
    zypper_install sudo curl python3 python3-pip ca-certificates
  else
    # best-effort generic
    if command -v apt-get >/dev/null 2>&1; then
      apt_install sudo curl python3 python3-venv python3-pip ca-certificates
    elif command -v dnf >/dev/null 2>&1 || command -v yum >/dev/null 2>&1; then
      yum_or_dnf_install sudo curl python3 python3-pip ca-certificates
    elif command -v apk >/dev/null 2>&1; then
      apk_install sudo curl python3 py3-pip
    else
      echo "Unsupported or unknown distro; please install sudo, curl, and python3 manually." >&2
      exit 2
    fi
  fi
}

install_docker() {
  if command -v docker >/dev/null 2>&1; then
    echo "Docker already installed: $(docker --version || true)"
    return
  fi

  if [[ $NONINTERACTIVE -eq 1 ]]; then
    export CHANNEL="stable"
  fi

  echo "Installing Docker using the official convenience script..."
  curl -fsSL https://get.docker.com -o /tmp/get-docker.sh
  sh /tmp/get-docker.sh

  if command -v systemctl >/dev/null 2>&1; then
    systemctl enable docker || true
    systemctl start docker || true
  elif command -v service >/dev/null 2>&1; then
    service docker start || true
  fi

  rm -f /tmp/get-docker.sh
  echo "Docker installed: $(docker --version || true)"
}

add_user_to_docker_group() {
  local u="$1"
  [[ -z "$u" ]] && return 0

  if ! id "$u" >/dev/null 2>&1; then
    echo "User '$u' does not exist; creating..."
    if command -v useradd >/dev/null 2>&1; then
      useradd -m -s /bin/bash "$u"
    elif command -v adduser >/dev/null 2>&1; then
      adduser -D "$u" || adduser "$u"
    else
      echo "No useradd/adduser found; cannot create user." >&2
      return 1
    fi
  fi

  if ! getent group docker >/dev/null 2>&1; then
    groupadd docker
  fi

  usermod -aG docker "$u"
  echo "User '$u' added to 'docker' group. User must re-login for group to take effect."
}

post_checks() {
  echo "Verifying Docker daemon..."
  if ! docker info >/dev/null 2>&1; then
    echo "WARNING: 'docker info' failed. The daemon may not be running or this shell lacks group membership. Trying to start..."
    if command -v systemctl >/dev/null 2>&1; then
      systemctl start docker || true
    fi
  fi

  echo "Ensuring Python available for Ansible..."
  if ! command -v python3 >/dev/null 2>&1; then
    echo "ERROR: Python3 not found after installation." >&2
    exit 3
  fi
}

main() {
  install_base_prereqs
  install_docker
  add_user_to_docker_group "${USER_TO_ADD}"
  post_checks

  echo
  echo "Bootstrap complete."
  echo "  - Docker: $(docker --version || echo 'not found')"
  echo "  - Python3: $(python3 --version 2>/dev/null || echo 'not found')"
  if [[ -n "${USER_TO_ADD}" ]]; then
    echo "  - Added '${USER_TO_ADD}' to 'docker' group (re-login required)."
  fi
}

main