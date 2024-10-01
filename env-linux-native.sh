#!/usr/bin/env bash
set -e

case "$(uname -m)" in
  "aarch64")
    export TARGET_ARCH="arm64"
    ;;
  "x86_64")
    export TARGET_ARCH="amd64"
    ;;
  *)
    printf "\e[1m\e[31m%s\e[0m\n" "Unsupported TARGET ARCH: '${TARGET_ARCH}'."
    exit 1
    ;;
esac

export GO="$(command -v go)"
