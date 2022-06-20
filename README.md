## devops
fast ops v2

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

### 镜像构建
```shell
#!/bin/sh
sh build.sh
```

### 依赖环境镜像(shell外壳)
[devops-shell](dockerfiles/devops-shell)
[devops-shell-java](dockerfiles/devops-shell-java)
[devops-shell-java-node](dockerfiles/devops-shell-java-node)

### 最终应用镜像
[devops](Dockerfile)
[devops-java](dockerfiles/devops-java)
[devops-java-node](dockerfiles/devops-java-node)
