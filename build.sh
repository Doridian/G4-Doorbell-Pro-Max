#!/bin/sh
set -ex

docker build -t g4adv-builder .
docker run --rm -it -v $(pwd):/src g4adv-builder /docker/build.sh
