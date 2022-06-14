#!/bin/sh

path=$(pwd)

#build shell
cd $path/dockerfiles/devops-shell
docker build --platform=linux/amd64 -t devops-shell .

cd $path/dockerfiles/devops-shell-java
docker build --platform=linux/amd64 -t devops-shell-java .

cd $path/dockerfiles/devops-shell-java-node
docker build --platform=linux/amd64 -t devops-shell-java-node .

#upload shell
docker tag devops-shell cubegroup/devops-shell
docker push cubegroup/devops-shell
docker tag devops-shell-java cubegroup/devops-shell-java
docker push cubegroup/devops-shell-java
docker tag devops-shell-java-node cubegroup/devops-shell-java-node
docker push cubegroup/devops-shell-java-node

#build app
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/app

#build devops
cd $path
docker build --platform=linux/amd64 -t devops .
docker tag devops cubegroup/devops
docker push cubegroup/devops

cd $path
docker build --platform=linux/amd64 -t devops-java -f $path/dockerfiles/devops-java/Dockerfile $path
docker tag devops-java cubegroup/devops-java
docker push cubegroup/devops-java

cd $path
docker build --platform=linux/amd64 -t devops-java-node -f $path/dockerfiles/devops-java-node/Dockerfile $path
docker tag devops-java-node cubegroup/devops-java-node
docker push cubegroup/devops-java-node
