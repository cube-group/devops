FROM --platform=linux/amd64 maven:3.8.6-eclipse-temurin-8-alpine
#devops程序运行所需外壳环境
USER root

COPY docker-entrypoint.sh /docker-entrypoint.sh

RUN echo "https://mirrors.aliyun.com/alpine/v3.15/main" > /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.15/community" >> /etc/apk/repositories && \
    apk update  && \
    apk add docker-cli openssh sshpass git tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata && \
    rm -rf /var/cache/apk

WORKDIR /
ENTRYPOINT ["sh","/docker-entrypoint.sh"]