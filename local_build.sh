#!/bin/sh
set -ex

mkdir -p ./dist/bin ./dist/lib

if [ -d /fw/image/rootfs ]; then
    rsync -av /fw/image/rootfs/lib/*.so* ./dist/lib
fi

REAL_CC=aarch64-linux-gnu-gcc
if [ "$(uname -m)" = "aarch64" ]; then
    REAL_CC=gcc
fi

export SRCDIR="$(pwd)"

export CC="$REAL_CC"
export CC_FOR_TARGET="$REAL_CC"

export CFLAGS="-g -O2 '-I${SRCDIR}/include'"
export LDFLAGS="$CFLAGS '-L${SRCDIR}/dist/lib' '${SRCDIR}/dist/lib/libc.so.6' '${SRCDIR}/dist/lib/ld-linux-aarch64.so.1' '${SRCDIR}/dist/lib/libpthread.so.0'"

runcc() {
    $CC_FOR_TARGET -Wall -Werror -nodefaultlibs $LDFLAGS "$@" -fno-stack-protector
}

export CGO_CFLAGS="$CFLAGS"
export CGO_LDFLAGS="$LDFLAGS"
export GOOS=linux
export GOARCH=arm64
export CGO_ENABLED=1
go mod tidy
go build -buildvcs=false -o ./dist/bin/g4adv ./cmd/g4adv
