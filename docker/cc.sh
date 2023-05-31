#!/bin/sh
set -e

REAL_CC=aarch64-linux-gnu-gcc
if [ "$(uname -s -m)" = "Linux aarch64" ]; then
    REAL_CC=gcc
fi

exec "$REAL_CC" "$@"
