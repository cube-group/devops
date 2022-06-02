#devops程序运行所需外壳环境
FROM --platform=linux/amd64 debian:11
USER root

COPY sources.list /etc/apt/sources.list
COPY docker.list /etc/apt/sources.list.d/docker.list
COPY gotty_linux_amd64.tar.gz /gotty_linux_amd64.tar.gz

RUN apt-get update -y && \
    apt-get install -y ca-certificates gpg dirmngr curl gpg-agent && \
    curl -fsSL "https://download.docker.com/linux/debian/gpg" | apt-key add -qq - && \
    apt-get update -y && \
    apt-get install -y tzdata ssh sshpass git docker-ce-cli --no-install-recommends && \
    cd / && tar -zxvf /gotty_linux_amd64.tar.gz >> /dev/null && \
    mv gotty /usr/local/bin/gotty && \
    chmod +x /usr/local/bin/gotty && \
    echo "Asia/Shanghai" > /etc/timezone && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    rm -rf /gotty_linux_amd64.tar.gz && \
    rm -rf /var/cache/* && \
    rm -rf /var/log/* && \
    rm -rf /var/lib/apt/lists/* && \
    apt autoclean && apt clean && apt autoremove

WORKDIR /