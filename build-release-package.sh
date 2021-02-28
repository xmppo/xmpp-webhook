#!/usr/bin/env bash

set -xe

git checkout "$1"
docker run --rm -ti  -v "$(pwd)":/build golang:1.15-buster sh -c "cd /build && go build"
tar -czvf "xmpp-webhook-$1-linux-amd64.tar.gz" xmpp-webhook xmpp-webhook.service README.md LICENSE THIRD-PARTY-NOTICES
sha512sum "xmpp-webhook-$1-linux-amd64.tar.gz" > "xmpp-webhook-$1-linux-amd64.tar.gz.sha512"
