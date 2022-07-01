FROM --platform=linux/amd64 alpine:3.15
#devops程序运行所需外壳环境
USER root

COPY entrypoint.sh /entrypoint.sh
COPY gotty /usr/local/bin/gotty

RUN echo "https://mirrors.aliyun.com/alpine/v3.15/main" > /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.15/community" >> /etc/apk/repositories && \
    apk update  && \
    apk add docker-cli openssh sshpass bash git tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    chmod +x /usr/local/bin/gotty && \
    apk del tzdata && \
    rm -rf /var/cache/apk

WORKDIR /
ENTRYPOINT ["sh","/entrypoint.sh"]