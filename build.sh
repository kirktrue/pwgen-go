#!/bin/bash -e

VERSION=$(cat version.txt)

function build() {
  local os=$1
  local arch=$2
  local file_name=pwgen-$VERSION-$os-$arch

  GOOS=$os GOARCH=$arch go build -o $file_name
}

build linux amd64
build darwin amd64
