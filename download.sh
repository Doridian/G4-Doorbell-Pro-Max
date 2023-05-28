#!/bin/sh
set -ex

# apt install jq curl wget python3-pip
# pip3 install --user ubi_reader

API_LATEST_URL='https://fw-update.ubnt.com/api/firmware-latest?filter=eq~~product~~uvc&filter=eq~~channel~~release&filter=eq~~platform~~s5l'

BIN_URL="$(curl -s -f "$API_LATEST_URL" | jq -r ._embedded.firmware[0]._links.data.href)"

BIN_NAME="fwdownload/firmware.bin"

if [ ! -f "$BIN_NAME" ]
then
    wget "$BIN_URL" -O "$BIN_NAME"
    rm -rf fwimage
    mkdir -p fwimage
    ubireader_extract_files "$BIN_NAME" -o fwimage
    mv fwimage/*/rootfs fwimage/rootfs
fi
