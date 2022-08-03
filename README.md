## devops
fast ops on docker

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
cubegroup/devops
```

### 应用镜像
* [devops](dockerfiles/devops)  对应[docker hub](https://hub.docker.com/r/cubegroup/devops)
* [devops-java](dockerfiles/devops-java) 对应[docker hub](https://hub.docker.com/r/cubegroup/devops-java)

### linux安装或升级docker
```shell
curl -fsSL https://get.docker.com/ | sh
```