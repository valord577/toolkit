#!/usr/bin/env bash
set -e

TARGET_PLATFORM=${1}
TARGET_ARCH=${2}

case ${TARGET_PLATFORM} in
  "macosx")
    TARGET_FLAG="macosx"
    ;;
  "iphoneos")
    TARGET_FLAG="iphoneos"
    ;;
  "iphonesimulator")
    TARGET_FLAG="ios-simulator"
    ;;
  *)
    printf "\e[1m\e[31m%s\e[0m\n" "Unsupported TARGET PLATFORM: '${TARGET_PLATFORM}'."
    exit 1
    ;;
esac

case ${TARGET_ARCH} in
  "arm64")
    export GOARCH="arm64"
    ;;
  "x86_64")
    export GOARCH="amd64"; export GOAMD64="v3"
    ;;
  *)
    printf "\e[1m\e[31m%s\e[0m\n" "Unsupported TARGET ARCH: '${TARGET_ARCH}'."
    exit 1
    ;;
esac

TARGET_DEPLOYMENT="10"
if [ "${TARGET_PLATFORM}" == "macosx" ]; then
  TARGET_DEPLOYMENT="10.15"
fi
SYSROOT="$(xcrun --sdk ${TARGET_PLATFORM} --show-sdk-path)"
CROSS_FLAGS="-arch ${TARGET_ARCH} -m${TARGET_FLAG}-version-min=${TARGET_DEPLOYMENT}"

export CGO_CFLAGS="$(go env CGO_CFLAGS) ${CROSS_FLAGS} --sysroot=${SYSROOT}"
export CGO_CXXFLAGS="$(go env CGO_CXXFLAGS) ${CROSS_FLAGS} --sysroot=${SYSROOT}"
export CGO_LDFLAGS="$(go env CGO_LDFLAGS) ${CROSS_FLAGS} --sysroot=${SYSROOT}"
export CGO_ENABLED="1"

export GO="$(command -v go)"
