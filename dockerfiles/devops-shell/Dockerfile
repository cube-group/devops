FROM --platform=linux/amd64 alpine:3.16
#devops程序运行所需外壳环境
USER root

COPY gotty /usr/local/bin/gotty

RUN echo "https://mirrors.aliyun.com/alpine/v3.16/main" > /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.16/community" >> /etc/apk/repositories && \
    apk update  && \
    apk add docker-cli openssh sshpass bash git tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    chmod +x /usr/local/bin/gotty && \
    apk del tzdata

WORKDIR /