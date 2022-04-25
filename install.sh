#!/bin/sh
# Usage: [sudo] [BINDIR=/usr/local/bin] ./install.sh [<BINDIR>]
#
# Examples:
#     # Install globally into /usr/local/bin
#     $ sudo ./install.sh /usr/local/bin
#
#     # Install globally into /usr/bin
#     $ sudo ./install.sh /usr/bin
#
#     # Install into a local directory for your user only
#     $ ./install.sh $HOME/usr/bin
#
#     # Install into a local directory for your user only, by defining
#     # the BINDIR environment variable
#     $ BINDIR=$HOME/usr/bin ./install.sh
#
# Default BINDIR=/usr/bin

# Note this install tool borrows heavily from the install tool found at:
#     https://github.com/zaquestion/lab/blob/58af7c16e737737a53fcd3430d91863e29f52129/install.sh

set -euf

if [ -n "${DEBUG-}" ]; then
    set -x
fi

: "${BINDIR:=/usr/bin}"

if [ $# -gt 0 ]; then
  BINDIR=$1
fi

_can_install() {
  if [ ! -d "${BINDIR}" ]; then
    mkdir -p "${BINDIR}" 2> /dev/null
  fi
  [ -d "${BINDIR}" ] && [ -w "${BINDIR}" ]
}

if ! _can_install && [ "$(id -u)" != 0 ]; then
  printf "Run script as sudo\n"
  exit 1
fi

if ! _can_install; then
  printf -- "Can't install to %s\n" "${BINDIR}"
  exit 1
fi

case "$(uname -m)" in
    x86_64)
        machine="x86_64"
        ;;
    i386)
        machine="i386"
        ;;
    *)
        machine=""
        ;;
esac

case $(uname -s) in
    Linux)
        os="Linux"
        ;;
    Darwin)
        os="Darwin"
        ;;
    *)
        printf "OS not supported by this installation script\n"
        printf "To see list of all builds, visit:\n"
        printf "    https://github.com/lelandbatey/histogram_timestamps/releases/tag/latest\n"
        exit 1
        ;;
esac

printf "Fetching latest version\n"
latest="$(curl -sL 'https://api.github.com/repos/lelandbatey/histogram_timestamps/releases/latest' | grep 'tag_name' | grep -o 'v[0-9\.]\+' | cut -c 2-)"
tempFolder="/tmp/histogram_timestamps_v${latest}"

printf -- "Found version %s\n" "${latest}"

mkdir -p "${tempFolder}" 2> /dev/null
printf -- "Downloading histogram_timestamps_%s_%s_%s.tar.gz\n" "${latest}" "${os}" "${machine}"
curl -sL "https://github.com/lelandbatey/histogram_timestamps/releases/download/v${latest}/histogram_timestamps_${latest}_${os}_${machine}.tar.gz" | tar -C "${tempFolder}" -xzf -

printf -- "Installing...\n"
install -m755 "${tempFolder}/histogram_timestamps" "${BINDIR}/histogram_timestamps"

printf "Cleaning up temp files\n"
rm -rf "${tempFolder}"

printf -- "Successfully installed histogram_timestamps in %s/\n" "${BINDIR}"
