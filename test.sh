#!/bin/sh
set -ex
rm -f mynfc myfp

# You need the following from the G4 Doorbell Pro firmware:
# - The entire /lib folder (copy it into ./rootfs/lib, such that ./rootfs/lib/libbmkt.so exists, etc)
# - Make sure your compiler uses ld-linux-aarch64.so.1 (found in /lib) from the rootfs (I just overworte /usr/aarch64-linux-gnu/lib/ld-linux-aarch64.so.1)

#aarch64-linux-gnu-gcc-12 -nodefaultlibs  ./mynfc.c -L./rootfs/lib -lnxp-nfc -o mynfc ./rootfs/lib/libc.so.6 -fno-stack-protector
aarch64-linux-gnu-gcc-12 -nodefaultlibs  ./myfp.c -L./rootfs/lib -lbmkt -lpthread-2.23 -o myfp ./rootfs/lib/libc.so.6 -fno-stack-protector
#cat mynfc | ssh ubnt@camera-front-door.foxden.network 'rm -f /var/mynfc && echo Loading... && cat /dev/stdin > /var/mynfc && echo CHMod && chmod 755 /var/mynfc'
cat myfp | ssh ubnt@camera-front-door.foxden.network 'rm -f /var/myfp && echo Loading... && cat /dev/stdin > /var/myfp && echo CHMod && chmod 755 /var/myfp'
