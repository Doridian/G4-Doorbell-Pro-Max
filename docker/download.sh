#!/bin/sh
set -ex

# apt install jq curl wget python3-pip
# pip3 install --user ubi_reader

#API_LATEST_URL='https://fw-update.ubnt.com/api/firmware-latest?filter=eq~~product~~uvc&filter=eq~~channel~~release&filter=eq~~platform~~s5l'
#BIN_URL="$(curl -s -f "$API_LATEST_URL" | jq -r ._embedded.firmware[0]._links.data.href)"

BIN_URL="https://fw-download.ubnt.com/data/uvc/f2ad-s5l-4.63.22-33aae7792c344d8db134fc91cbb99290.bin"

mkdir -p /fw/download
# TODO: Use unifi's filename
BIN_NAME="/fw/download/firmware.bin"

if [ ! -f "$BIN_NAME" ]
then
    wget "$BIN_URL" -O "$BIN_NAME"
    rm -rf /fw/image
    mkdir -p /fw/image
    ubireader_extract_files "$BIN_NAME" -o /fw/image
    mv /fw/image/*/rootfs /fw/image/rootfs
    rmdir /fw/image/* || true
fi
