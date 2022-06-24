FROM alpine:3.15
#devops程序运行所需外壳环境
USER root

COPY entrypoint.sh /entrypoint.sh

RUN echo "https://mirrors.aliyun.com/alpine/v3.15/main" > /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.15/community" >> /etc/apk/repositories && \
    apk update  && \
    apk add openssh sshpass tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata && \
    rm -rf /var/cache/apk

WORKDIR /

ENTRYPOINT ["sh","/entrypoint.sh"]