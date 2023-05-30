#!/bin/sh
set -ex

mkdir -p ./dist/bin ./dist/lib

#cp /fw/image/rootfs/lib/*.so* ./dist/lib

REAL_CC=aarch64-linux-gnu-gcc
if [ "$(uname -m)" = "aarch64" ]; then
    REAL_CC=gcc
fi

export CC="$REAL_CC"
export CC_FOR_TARGET="$REAL_CC"

export CFLAGS="-g -O2 -I/src/include"
export LDFLAGS="$CFLAGS -L/fw/image/rootfs/lib -lbmkt -lnxp-nfc /fw/image/rootfs/lib/libc.so.6 /fw/image/rootfs/lib/ld-linux-aarch64.so.1 /fw/image/rootfs/lib/libpthread.so.0"

runcc() {
    $CC_FOR_TARGET -Wall -Werror -nodefaultlibs $LDFLAGS "$@" -fno-stack-protector
}

export CGO_CFLAGS="$CFLAGS"
export CGO_LDFLAGS="$LDFLAGS"
export GOOS=linux
export GOARCH=arm64
export CGO_ENABLED=1
go mod tidy
go build -o ./dist/bin/g4adv ./cmd/g4adv
