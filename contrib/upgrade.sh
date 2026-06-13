#!/usr/bin/env bash
# This is an update script for CreeperCoding installed via the binary distribution
# on linux as systemd service. It performs a backup and updates
# CreeperCoding in place.
# NOTE: This adds the GPG Signing Key of CreeperCoding maintainers to the keyring.
# Depends on: bash, curl, xz, sha256sum. optionally jq, gpg
#   See section below for available environment vars.
#   When no version is specified, updates to the latest release.
# Examples:
#   upgrade.sh 1.15.10
  # creepercodinghome=/opt/creepercoding creepercodingconf=$creepercodinghome/app.ini upgrade.sh

# Check if CreeperCoding service is running
if ! pidof creepercoding &> /dev/null; then
  echo "Error: CreeperCoding is not running."
  exit 1
fi

# Continue with rest of the script if CreeperCoding is running
echo "CreeperCoding is running. Continuing with rest of script..."

# apply variables from environment
: "${creepercodingbin:="/usr/local/bin/creepercoding"}"
: "${creepercodinghome:="/var/lib/creepercoding"}"
: "${creepercodingconf:="/etc/creepercoding/app.ini"}"
: "${creepercodinguser:="git"}"
: "${sudocmd:="sudo"}"
: "${arch:="linux-amd64"}"
: "${service_start:="$sudocmd systemctl start creepercoding"}"
: "${service_stop:="$sudocmd systemctl stop creepercoding"}"
: "${service_status:="$sudocmd systemctl status creepercoding"}"
: "${backupopts:=""}" # see `creepercoding dump --help` for available options

function creepercodingcmd {
  if [[ $sudocmd = "su" ]]; then
    # `-c` only accept one string as argument.
    "$sudocmd" - "$creepercodinguser" -c "$(printf "%q " "$creepercodingbin" "--config" "$creepercodingconf" "--work-path" "$creepercodinghome" "$@")"
  else
    "$sudocmd" --user "$creepercodinguser" "$creepercodingbin" --config "$creepercodingconf" --work-path "$creepercodinghome" "$@"
  fi
}

function require {
  for exe in "$@"; do
    command -v "$exe" &>/dev/null || (echo "missing dependency '$exe'"; exit 1)
  done
}

# parse command line arguments
while true; do
  case "$1" in
    -v | --version ) creepercodingversion="$2"; shift 2 ;;
    -y | --yes ) no_confirm="yes"; shift ;;
    --ignore-gpg) ignore_gpg="yes"; shift ;;
    "" | -- ) shift; break ;;
    * ) echo "Usage:  [<environment vars>] upgrade.sh [-v <version>] [-y] [--ignore-gpg]"; exit 1;; 
  esac
done

# exit once any command fails. this means that each step should be idempotent!
set -euo pipefail

if [[ -f /etc/os-release ]]; then
  os_release=$(cat /etc/os-release)

  if [[ "$os_release" =~ "OpenWrt" ]]; then
    sudocmd="su"
   service_start="/etc/init.d/creepercoding start"
   service_stop="/etc/init.d/creepercoding stop"
   service_status="/etc/init.d/creepercoding status"
  else
    require systemctl
  fi
fi

require curl xz sha256sum "$sudocmd"

# select version to install
if [[ -z "${creepercodingversion:-}" ]]; then
  require jq
  creepercodingversion=$(curl --connect-timeout 10 -sL https://dl.gitea.com/gitea/version.json | jq -r .latest.version)
  echo "Latest available version is $creepercodingversion"
fi

# confirm update
echo "Checking currently installed version..."
current=$(creepercodingcmd --version | cut -d ' ' -f 3)
[[ "$current" == "$creepercodingversion" ]] && echo "$current is already installed, stopping." && exit 0
if [[ -z "${no_confirm:-}"  ]]; then
  echo "Make sure to read the changelog first: https://github.com/go-gitea/gitea/blob/main/CHANGELOG.md"
  echo "Are you ready to update CreeperCoding from ${current} to ${creepercodingversion}? (y/N)"
  read -r confirm
  [[ "$confirm" == "y" ]] || [[ "$confirm" == "Y" ]] || exit 1
fi

echo "Upgrading CreeperCoding from $current to $creepercodingversion ..."

pushd "$(pwd)" &>/dev/null
cd "$creepercodinghome" # needed for creepercoding dump later

# download new binary
binname="creepercoding-${creepercodingversion}-${arch}"
binurl="https://dl.gitea.com/gitea/${creepercodingversion}/${binname}.xz"
echo "Downloading $binurl..."
curl --connect-timeout 10 --silent --show-error --fail --location -O "$binurl{,.sha256,.asc}"

# validate checksum & gpg signature
sha256sum -c "${binname}.xz.sha256"
if [[ -z "${ignore_gpg:-}" ]]; then
  require gpg
  # try to use curl first, it uses standard tcp 443 port and works better behind strict firewall rules
  curl -fsSL --connect-timeout 10 "https://keys.openpgp.org/vks/v1/by-fingerprint/7C9E68152594688862D62AF62D9AE806EC1592E2" | gpg --import \
    || gpg --keyserver keys.openpgp.org --recv 7C9E68152594688862D62AF62D9AE806EC1592E2
  gpg --verify "${binname}.xz.asc" "${binname}.xz" || { echo 'Signature does not match'; exit 1; }
fi
rm "${binname}".xz.{sha256,asc}

# unpack binary + make executable
xz --decompress --force "${binname}.xz"
chown "$creepercodinguser" "$binname"
chmod +x "$binname"

# stop CreeperCoding, create backup, replace binary, restart CreeperCoding
echo "Flushing CreeperCoding queues at $(date)"
creepercodingcmd manager flush-queues
echo "Stopping CreeperCoding at $(date)"
$service_stop
echo "Creating backup in $creepercodinghome"
# shellcheck disable=SC2086 # flag string
creepercodingcmd dump $backupopts
echo "Updating binary at $creepercodingbin"
cp -f "$creepercodingbin" "$creepercodingbin.bak" && mv -f "$binname" "$creepercodingbin"
# Restore SELinux context if applicable (e.g. RHEL/Fedora)
command -v restorecon &>/dev/null && restorecon -v "$creepercodingbin" || true
$service_start
$service_status

echo "Upgrade to $creepercodingversion successful!"

popd
