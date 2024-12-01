#!/bin/bash

sudo apt update
sudo apt install --no-install-recommends -y \
    swtpm swtpm-tools

PREFIX="/usr/local" && \
VERSION="1.42.0" && \
curl -sSL \
"https://github.com/bufbuild/buf/releases/download/v${VERSION}/buf-$(uname -s)-$(uname -m).tar.gz" | \
sudo tar -xvzf - -C "${PREFIX}" --strip-components 1

go install github.com/goreleaser/goreleaser/v2@latest
