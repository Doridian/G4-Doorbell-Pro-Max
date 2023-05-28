#!/bin/sh
set -ex
rm -f mynfc myfp

# NOTE: Make sure your compiler uses ld-2.23.so (found in /lib) from the rootfs (I just overworte /usr/aarch64-linux-gnu/lib/ld-linux-aarch64.so.1)

#aarch64-linux-gnu-gcc-12 -nodefaultlibs  ./mynfc.c -L./fwimage/rootfs/lib -lnxp-nfc -o mynfc ./fwimage/rootfs/lib/libc.so.6 -fno-stack-protector
aarch64-linux-gnu-gcc-12 -nodefaultlibs  ./myfp.c -L./fwimage/rootfs/lib -lbmkt -lpthread-2.23 -o myfp ./fwimage/rootfs/lib/libc.so.6 -fno-stack-protector
#cat mynfc | ssh ubnt@camera-front-door.foxden.network 'rm -f /var/mynfc && echo Loading... && cat /dev/stdin > /var/mynfc && echo CHMod && chmod 755 /var/mynfc'
cat myfp | ssh ubnt@camera-front-door.foxden.network 'rm -f /var/myfp && echo Loading... && cat /dev/stdin > /var/myfp && echo CHMod && chmod 755 /var/myfp'
