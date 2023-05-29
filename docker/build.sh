#!/bin/sh
set -ex

CC=aarch64-linux-gnu-gcc
if [ "$(uname -m)" = "aarch64" ]; then
    CC=gcc
fi

runcc() {
    $CC -O3 -nodefaultlibs -L/fw/image/rootfs/lib "$@" /fw/image/rootfs/lib/libc.so.6 /fw/image/rootfs/lib/ld-2.23.so -fno-stack-protector
}

mkdir -p ./dist/bin ./dist/lib

cp /fw/image/rootfs/lib/libnxp-nfc.so /fw/image/rootfs/lib/libbmkt.so ./dist/lib
runcc mynfc.c -o ./dist/bin/mynfc -lnxp-nfc
runcc myfp.c -o ./dist/bin/myfp -lbmkt -lpthread-2.23
