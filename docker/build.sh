#!/bin/sh
set -ex

mkdir -p ./dist/bin ./dist/lib

cp /fw/image/rootfs/lib/libnxp-nfc.so /fw/image/rootfs/lib/libbmkt.so ./dist/lib

REAL_CC=aarch64-linux-gnu-gcc
if [ "$(uname -m)" = "aarch64" ]; then
    REAL_CC=gcc
fi

#export CC=/docker/cc.sh
#export CC_FOR_TARGET=/docker/cc.sh
export CC="$REAL_CC"
export CC_FOR_TARGET="$REAL_CC"

export LDFLAGS="-g -O2 -L/fw/image/rootfs/lib /fw/image/rootfs/lib/libc.so.6 /fw/image/rootfs/lib/ld-linux-aarch64.so.1 /fw/image/rootfs/lib/libpthread.so.0"

runcc() {
    $CC_FOR_TARGET -Wall -Werror -nodefaultlibs $LDFLAGS "$@" -fno-stack-protector
}
#runcc mynfc.c -o ./dist/bin/mynfc -lnxp-nfc
#runcc myfp.c -o ./dist/bin/myfp -lbmkt

export CGO_LDFLAGS="$LDFLAGS"
export GOOS=linux
export GOARCH=arm64
export CGO_ENABLED=1
go build -o ./dist/bin/g4adv ./cmd/g4adv
