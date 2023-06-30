FROM ubuntu:22.04

ENV LC_ALL C.UTF-8
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US.UTF-8

RUN apt update && \
    apt install -y skopeo wget &&\
    wget -O image-sync_1.0.2_linux_amd64.tar.gz https://github.com/cnsync/image-sync/releases/download/v1.0.2/image-sync_1.0.2_linux_amd64.tar.gz && \
    tar -zxvf image-sync_1.0.2_linux_amd64.tar.gz -C /usr/bin/ &&\
    chmod +x /usr/bin/image-sync

ADD entrypoint.sh /
RUN chmod +x /entrypoint.sh

WORKDIR /github/workspace
ENTRYPOINT ["/entrypoint.sh"]