#!/bin/sh
set -ex

CC=aarch64-linux-gnu-gcc
if [ "$(uname -m)" = "aarch64" ]; then
    CC=gcc
fi

runcc() {
    $CC -nodefaultlibs -L/fw/image/rootfs/lib "$@" /fw/image/rootfs/lib/libc.so.6 /fw/image/rootfs/lib/ld-2.23.so -fno-stack-protector
}

copycam() {
    # TODO: scp
    cat "$1" | ssh ubnt@camera-front-door.foxden.network "rm -f '/var/$1' && echo Loading... && cat /dev/stdin > '/var/$1' && echo CHMod && chmod 755 '/var/$1'"
}

runcc ./mynfc.c -o mynfc -lnxp-nfc
runcc ./myfp.c -o myfp -lbmkt -lpthread-2.23
