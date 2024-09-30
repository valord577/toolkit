#!/usr/bin/env bash
set -e

# ----------------------------
# Release Build
# ----------------------------
GO_BUILD_GCFLAGS=""
GO_BUILD_LDFLAGS="-s -w"
GO_MACRO_FLAVOR="-X '${PKG_NAME}/system.flavor=release'"
# ----------------------------
# Software Version
# ----------------------------
if command -v git >/dev/null 2>&1 ; then
  pushd "${PROJ_ROOT}"
  GIT_HASH="$(git describe --tags --always --dirty --abbrev=${GIT_ABBREV:-"7"})"
  popd
else
  printf "\e[1m\e[33m%s\e[0m\n" "@@@ Warn - unknown command: git"
  GIT_HASH=${GIT_HASH:-"unknown"}
fi
GO_BUILD_VERSION="${GIT_HASH}"
GO_MACRO_VERSION="-X '${PKG_NAME}/system.version=${GO_BUILD_VERSION}'"
# ----------------------------
# Time to be compiled
# ----------------------------
BUILD_DATE=$(date -u '+%Y-%m-%dT%H:%M:%SZ%:z')
GO_MACRO_DATETIME="-X '${PKG_NAME}/system.datetime=${BUILD_DATE}'"
# ----------------------------
# Start compiling :p
# ----------------------------
export GO111MODULE="on"

GO_BUILD_COMMAND=$(cat <<- EOF
${GO} build -o '${PKG_INST_DIR}/${PKG_NAME}$(go env GOEXE)' ${GO_BUILD_GCFLAGS} \
  -ldflags="-v ${GO_BUILD_LDFLAGS} ${GO_MACRO_VERSION} ${GO_MACRO_DATETIME} ${GO_MACRO_FLAVOR}" \
  '${PROJ_ROOT}'
EOF
)
printf "\e[1m\e[36m%s\e[0m\n" "${GO_BUILD_COMMAND}"; ${GO} env
eval ${GO_BUILD_COMMAND} \
  || { ret=$?; printf "\e[1m\e[31m%s\e[0m\n" "Failed to build golang exec: '${PKG_NAME}'."; exit "$ret"; }
# ----------------------------
# Print information
# ----------------------------
if command -v tree >/dev/null 2>&1 ; then
  tree -L 3 ${PKG_INST_DIR}
else
  ls -alh -- ${PKG_INST_DIR}
fi
printf "\e[1m\e[35m%s\e[0m\n" "${PKG_INST_DIR} - Build Done @${BUILD_DATE}"
