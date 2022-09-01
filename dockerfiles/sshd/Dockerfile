FROM --platform=linux/amd64 alpine:3.15

ENV ROOT_PASSWORD "root"

COPY run.sh /run.sh

RUN echo "https://mirrors.aliyun.com/alpine/v3.15/main" > /etc/apk/repositories && \
    echo "https://mirrors.aliyun.com/alpine/v3.15/community" >> /etc/apk/repositories && \
    apk update && \
    apk add docker-cli openssh-server tzdata openssh && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    sed -i "s/#PermitRootLogin.*/PermitRootLogin yes/g" /etc/ssh/sshd_config && \
    ssh-keygen -t rsa -P "" -f /etc/ssh/ssh_host_rsa_key && \
    ssh-keygen -t ecdsa -P "" -f /etc/ssh/ssh_host_ecdsa_key && \
    ssh-keygen -t ed25519 -P "" -f /etc/ssh/ssh_host_ed25519_key

# 开放22端口
EXPOSE 22
VOLUME /var/run/docker.sock
# 执行ssh启动命令
CMD ["sh", "/run.sh"]
