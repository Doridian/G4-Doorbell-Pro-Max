FROM ubuntu:24.04

RUN apt update && \
    apt install -y \
        git \
        curl \
        wget \
        jq \
        python3-pip \
        golang-go \
        rsync \
    && pip install ubi_reader

RUN [ "$(uname -m)" = "aarch64" ] && (apt install -y gcc && ln -s "$(which gcc)" /bin/gcc-aarch64-linux-gnu) || apt install -y gcc-aarch64-linux-gnu

WORKDIR /fw
COPY /docker/download.sh /docker/download.sh
RUN /docker/download.sh
COPY docker /docker

RUN mkdir -p /src
WORKDIR /src
VOLUME /src
