## devops
fast ops v2

### 运行
可使用hub.docker.com官方镜像
```shell
#!/bin/sh
docker run -it -d --restart=always \
--net=host \
--privileged \
-e WEB_PORT=8080 \
-e DB_HOST=xx \
-e DB_USERNAME=xx \
-e DB_PASSWORD=xx \
-e DB_NAME=xx \
-e WEB_DEBUG=0 \
#for git
-v ~/.ssh:/root/.ssh \ 
--name devops \
cubegroup/devops:v2
```

### 镜像构建
```shell
#!/bin/sh
sh build.sh
```