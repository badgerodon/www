#!/bin/bash

set -e

mkdir -p dist

docker build . -f build.dockerfile -t badgerodon-www-build
docker create --name extract badgerodon-www-build:latest
docker cp extract:/tmp/badgerodon-www.tar.xz dist/badgerodon-www.tar.xz
docker rm -f extract
