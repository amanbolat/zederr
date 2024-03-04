#!/bin/bash

# This script downloads and installs the protoc binary for the current platform.
# Only Linux and macOS are supported.
#
# Required environment variables:
# - PROTOC_VER: the version of protoc to install (e.g. 25.3)
# - PROTOC_CHECKSUM: the sha256 checksum of the protoc binary (e.g. d0fcd6d3b3ef6f22f1c47cc30a80c06727e1eccdddcaf0f4a3be47c070ffd3fe)
#
# To get the checksum run the command below:
# curl -L "https://github.com/protocolbuffers/protobuf/releases/download/v25.3/protoc-25.3-osx-aarch_64.zip" | sha256sum
set -ex

if [[ -z "$PROTOC_VER" ]]; then
  echo "PROTOC_VER environment variable not set"
  exit 1
fi

if [[ -z "$PROTOC_CHECKSUM" ]]; then
  echo "PROTOC_CHECKSUM environment variable not set"
  exit 1
fi

os_name="$(uname -o)"
case $os_name in
  "GNU/Linux")
    os_name="linux"
    ;;
  "Darwin")
    os_name="osx"
    ;;
  *)
    echo "unsupported OS: $os_name"
    exit 1
    ;;
esac

os_arch="$(uname -m)"
case $os_arch in
  "x86_64")
    os_arch="x86_64"
    ;;
  "arm64")
    os_arch="aarch_64"
    ;;
  *)
    echo "unsupported architecture: $os_arch"
    exit 1
    ;;
esac

if [ -f ".bin/protoc.zip" ]; then
  ACTUAL_CHECKSUM=$(sha256sum .bin/protoc.zip | awk '{ print $1 }')
  if [ "$ACTUAL_CHECKSUM" == "$PROTOC_CHECKSUM" ]; then
    echo "Checksums match, no need to download"
    exit 0
  fi
fi

curl -sSL \
  "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VER}/protoc-${PROTOC_VER}-${os_name}-${os_arch}.zip" \
  -o .bin/protoc.zip
rm -rf .bin/protocd
rm -rf .bin/protoc
mkdir -p .bin/protocd
unzip .bin/protoc.zip -d .bin/protocd
chmod +x .bin/protocd/bin/protoc
