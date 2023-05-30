#!/bin/sh
set -ex

REAL_CC=aarch64-linux-gnu-gcc
if [ "$(uname -m)" = "aarch64" ]; then
    REAL_CC=gcc
fi

if command -v "$REAL_CC" >/dev/null 2>/dev/null; then
    echo 'Found required GCC locally, using it...'
    ./local_build.sh
else
    echo 'Could not find required GCC, using Docker...'
    docker build -t g4adv-builder .
    docker run --rm -it -v $(pwd):/src g4adv-builder ./local_build.sh
fi

copycam() {
    cat "dist/bin/$1" | ssh ubnt@camera-front-door.foxden.network "rm -f '/var/$1' && echo Loading... && cat /dev/stdin > '/var/$1' && echo CHMod && chmod 755 '/var/$1'"
}

copycam g4adv
