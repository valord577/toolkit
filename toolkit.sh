#!/usr/bin/env bash
set -e

GO_MODULE="toolkit"

# ----------------------------
# Release Build
# ----------------------------
GO_BUILD_FLAVOR="release"
GO_BUILD_GCFLAGS=""
GO_BUILD_LDFLAGS="-s -w"

GO_MACRO_FLAVOR="-X '${GO_MODULE}/system.flavor=${GO_BUILD_FLAVOR}'"
# ----------------------------
# Supported OS/ARCH
# ----------------------------
SUPPORTED_TARGET=$(cat <<- 'EOF'
linux/amd64
linux/arm64
windows/amd64
darwin/amd64
darwin/arm64
EOF
)

GOOS=${GOOS:-"$(go env GOOS)"}
GOARCH=${GOARCH:-"$(go env GOARCH)"}
printf "\e[1m\e[35m%s\e[0m\n" "BUILD TARGET: '${GOOS}/${GOARCH}'."

UNSUPPORTED_ERR="1"
for t in ${SUPPORTED_TARGET[@]}; do
  if [ "${t}" == "${GOOS}/${GOARCH}" ]; then
    UNSUPPORTED_ERR="0"
  fi
done
if [ "${UNSUPPORTED_ERR}" == "1" ]; then
  printf "\e[1m\e[31m%s\e[0m\n" "Invalid BUILD TARGET: '${GOOS}/${GOARCH}'."
  exit 1
fi
# ----------------------------
# Directory of compiled outputs
# ----------------------------
PROJ_ROOT=$(cd "$(dirname ${BASH_SOURCE[0]})"; pwd)
GO_INST_DIR=${GO_INST_DIR:-"${PROJ_ROOT}/out_go/${GOOS}/${GOARCH}"}
# ----------------------------
# Software Version
# ----------------------------
if command -v git >/dev/null 2>&1 ; then
  pushd "${PROJ_ROOT}"
  GIT_HASH="$(git describe --tags --always --dirty --abbrev=${GIT_ABBREV:-"7"})"
  popd
else
  printf "\e[1m\e[33m%s\e[0m\n" "@@@ Warn - unknown command: git"
  GIT_HASH="unknown"
fi
GO_BUILD_VERSION="${GIT_HASH}"
GO_MACRO_VERSION="-X '${GO_MODULE}/system.version=${GO_BUILD_VERSION}'"
# ----------------------------
# Time to be compiled
# ----------------------------
BUILD_DATE=$(date -u '+%Y-%m-%dT%H:%M:%SZ%:z')
GO_MACRO_DATETIME="-X '${GO_MODULE}/system.datetime=${BUILD_DATE}'"
# ----------------------------
# Set ENV
# ----------------------------
export GO111MODULE="on"
export CGO_ENABLED="1"
export GOOS="${GOOS}"
export GOARCH="${GOARCH}"
export GOPROXY="https://goproxy.cn,direct"
export GOSUMDB="sum.golang.google.cn"
# ----------------------------
# Start compiling :p
# ----------------------------
GO_PROGRAMS=${GO_PROGRAMS:-""}
if [ -z "${GO_PROGRAMS}" ]; then
  GO_PROGRAMS="autoip autossh netdev"
fi

for prog in ${GO_PROGRAMS[@]}; do
  GO_BUILD_COMMAND=$(cat <<- EOF
go build -o '${GO_INST_DIR}/${prog}$(go env GOEXE)' ${GO_BUILD_GCFLAGS} \
  -ldflags="${GO_BUILD_LDFLAGS} ${GO_MACRO_VERSION} ${GO_MACRO_DATETIME} ${GO_MACRO_FLAVOR}" \
  '${PROJ_ROOT}/toolkit/${prog}'
EOF
)
  printf "\e[1m\e[36m%s\e[0m\n" "${GO_BUILD_COMMAND}"
  eval ${GO_BUILD_COMMAND} \
    || { ret=$?; printf "\e[1m\e[31m%s\e[0m\n" "Failed to build golang exec: '${prog}'."; exit "$ret"; }
done
# ----------------------------
# Print information
# ----------------------------
if command -v tree >/dev/null 2>&1 ; then
  tree ${GO_INST_DIR}
else
  ls -alh -- ${GO_INST_DIR}
fi
printf "\e[1m\e[35m%s\e[0m\n" "${GO_INST_DIR} - Build Done @${BUILD_DATE}"
