#!/usr/bin/env bash
set -e

TARGET_ARCH=${1}
TARGET_LIBC=${2}

case ${TARGET_ARCH} in
  "arm64")
    __TARGET_TRIPLE__="aarch64-unknown-linux-${TARGET_LIBC}"
    ;;
  "amd64")
    __TARGET_TRIPLE__="x86_64-pc-linux-${TARGET_LIBC}"
    ;;
  *)
    printf "\e[1m\e[31m%s\e[0m\n" "Unsupported TARGET ARCH: '${TARGET_ARCH}'."
    exit 1
    ;;
esac

CROSS_TOOLCHAIN_ROOT=${CROSS_TOOLCHAIN_ROOT:-""}
if [ -z "${CROSS_TOOLCHAIN_ROOT}" ]; then
  printf "\e[1m\e[31m%s\e[0m\n" "Blank CROSS_TOOLCHAIN_ROOT: '${CROSS_TOOLCHAIN_ROOT}'."
  exit 1
fi

pushd ${PROJ_ROOT}/cross/linux; { ln -sfn "pkgconf-wrapper" "pkgconf-wrapper.${__TARGET_TRIPLE__}"; }; popd
BUILTIN_CROSS_TOOLCHAIN_PKGCONF="${PROJ_ROOT}/cross/linux/pkgconf-wrapper.${__TARGET_TRIPLE__}"
if [ -n "${CROSS_TOOLCHAIN_PKGCONF_PREFIX}" ]; then
  export CROSS_TOOLCHAIN_PKGCONF="${CROSS_TOOLCHAIN_PKGCONF_PREFIX}.${__TARGET_TRIPLE__}"
else
  export CROSS_TOOLCHAIN_PKGCONF="${BUILTIN_CROSS_TOOLCHAIN_PKGCONF}"
fi

pushd ${PROJ_ROOT}/cross/linux; { ln -sfn "cgotool-wrapper" "cgotool-wrapper.${__TARGET_TRIPLE__}"; }; popd
BUILTIN_CROSS_TOOLCHAIN_CGOTOOL="${PROJ_ROOT}/cross/linux/cgotool-wrapper.${__TARGET_TRIPLE__}"
if [ -n "${CROSS_TOOLCHAIN_CGOTOOL_PREFIX}" ]; then
  export CROSS_TOOLCHAIN_CGOTOOL="${CROSS_TOOLCHAIN_CGOTOOL_PREFIX}.${__TARGET_TRIPLE__}"
else
  export CROSS_TOOLCHAIN_CGOTOOL="${BUILTIN_CROSS_TOOLCHAIN_CGOTOOL}"
fi
export GO="${CROSS_TOOLCHAIN_CGOTOOL}"
