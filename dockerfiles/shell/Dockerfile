FROM docker:20.10.16-dind-alpine3.15
#devops程序运行所需外壳环境
USER root

COPY gotty /usr/local/bin/gotty
COPY entrypoint.sh /entrypoint.sh

RUN echo "https://mirrors.aliyun.com/alpine/v3.15/main" > /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.15/community" >> /etc/apk/repositories && \
    apk update && \
    apk add sshpass bash tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    chmod +x /usr/local/bin/gotty && \
    apk del tzdata && \
    rm -rf /var/cache/apk

WORKDIR /go

ENTRYPOINT ["sh","-e","/entrypoint.sh"]