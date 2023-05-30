FROM ubuntu:22.04

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

RUN [ "$(uname -m)" = "aarch64" ] && apt install -y gcc || apt install -y gcc-aarch64-linux-gnu

WORKDIR /fw
COPY /docker/download.sh /docker/download.sh
RUN /docker/download.sh
COPY docker /docker

RUN mkdir -p /src /home/user && groupadd -g 1000 user && useradd -u 1000 -g 1000 user
WORKDIR /src
VOLUME /src
RUN chown 1000:1000 /src /home/user
USER 1000:1000
