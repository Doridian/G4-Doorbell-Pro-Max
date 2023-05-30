#!/bin/sh
set -ex

docker build -t g4adv-builder .
docker run --rm -it -v $(pwd):/src g4adv-builder /docker/build.sh

copycam() {
    cat "dist/bin/$1" | ssh ubnt@camera-front-door.foxden.network "rm -f '/var/$1' && echo Loading... && cat /dev/stdin > '/var/$1' && echo CHMod && chmod 755 '/var/$1'"
}

copycam g4adv
