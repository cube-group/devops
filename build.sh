#!/bin/sh

path=$(pwd)

#build shell
docker build --platform=linux/amd64 -t devops-shell -f $path/dockerfiles/devops-shell/Dockerfile $path/dockerfiles

docker build --platform=linux/amd64 -t devops-shell-java -f $path/dockerfiles/devops-shell-java/Dockerfile $path/dockerfiles

#upload shell
docker tag devops-shell cubegroup/devops-shell
docker push cubegroup/devops-shell
docker tag devops-shell-java cubegroup/devops-shell-java
docker push cubegroup/devops-shell-java

#build app
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/app

#build devops
docker build --platform=linux/amd64 -t devops $path
docker tag devops cubegroup/devops
docker push cubegroup/devops

docker build --platform=linux/amd64 -t devops-java -f $path/dockerfiles/devops-java/Dockerfile $path
docker tag devops-java cubegroup/devops-java
docker push cubegroup/devops-java