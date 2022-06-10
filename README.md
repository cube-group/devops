## devops
fast ops

### 运行
可使用hub.docker.com官方镜像
```shell
#!/bin/sh
docker run -it -d --restart=always \
-p 8080:80 \
-e DB_HOST=xx \
-e DB_USERNAME=xx \
-e DB_PASSWORD=xx \
-e DB_NAME=xx \
-e WEB_DEBUG=0 \
# for git
-v ~/.ssh:/root/.ssh \ 
-v /var/run/docker.sock:/var/run/docker.sock \
--name devops \
cube-group/devops
```

### 镜像构建
```shell
#!/bin/sh
sh build.sh
```
### 依赖环境镜像(shell外壳)
[Dockerfile](shell/Dockerfile)

### 应用镜像
[Dockerfile](Dockerfile)

### 添加其他程序
安装openjdk
```Dockerfile
FROM devops
#FROM cubegroup/devops
RUN apt update && \
    apt install -y apt-transport-https ca-certificates wget dirmngr gnupg software-properties-common && \
    wget -qO - https://adoptopenjdk.jfrog.io/adoptopenjdk/api/gpg/key/public | apt-key add - && \
    add-apt-repository --yes https://adoptopenjdk.jfrog.io/adoptopenjdk/deb/ && \
    apt update && \
    apt install -y adoptopenjdk-8-hotspot
```

安装golang 
```Dockerfile
FROM devops
#FROM cubegroup/devops
VOLUME /data
#golang
RUN wget https://studygolang.com/dl/golang/go1.16.7.linux-amd64.tar.gz && \
    tar -zxf go1.16.7.linux-amd64.tar.gz && \
    rm /go1.16.7.linux-amd64.tar.gz && \
    mv /go /usr/local/ && \
    ln -s /usr/local/go/bin/go /usr/local/bin/go && \
    mkdir -p /data/go && \
    go env -w GOPRIVATE=gitlab.xx.com && \
    git config --global url."git@gitlab.xx.com:".insteadOf "https://gitlab.xx.com/"
```

安装nodejs
```Dockerfile
FROM devops
#FROM cubegroup/devops
VOLUME /data
#nodejs
RUN apt install -y nodejs npm && \
    mkdir -p /data/nodejs/node_modules/node_global && \
    mkdir -p /data/nodejs/node_modules/node_cache && \
    npm config set unsafe-perm true && \
    npm config set registry https://registry.npm.taobao.org && \
    npm install -g npm && \
    npm config set prefix "/data/nodejs/node_modules/node_global" && \
    npm config set cache "/data/nodejs/node_modules/node_cache" && \
    #npm-cli-adduser
    npm install npm-cli-adduser -g && \
    ln -s /data/nodejs/node_modules/node_global/lib/node_modules/npm-cli-adduser/index.js /usr/local/bin/npm-cli-adduser && \
    #cnpm
    npm install -g cnpm --registry=https://registry.npmmirror.com
```

安装python3
```Dockerfile
FROM devops
#FROM cubegroup/devops
VOLUME /data
VOLUME /usr/lib/python3.7/site-packages
RUN apt install -y python3 python3-pip&& \
    mkdir -p /data/pypi/cache && \
    mkdir -p /usr/lib/python3.7/site-packages && \
    chmod -R 777 /usr/lib/python3.7/site-packages && \
    pip3 install --upgrade pip
```