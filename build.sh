#!/bin/sh

docker rmi -f devops devops-shell devops-java devops-java-node

path=$(pwd)
#image devops-shell
cd $path/shell
docker build --platform=linux/amd64 -t devops-shell .
docker tag devops-shell cubegroup/devops-shell
docker push cubegroup/devops-shell

#image devops
cd $path
docker build --platform=linux/amd64 -t devops .
docker tag devops cubegroup/devops
docker push cubegroup/devops

#image devops-java
cd $path/shell/java
docker build --platform=linux/amd64 -t devops-java .
docker tag devops-java cubegroup/devops-java
docker push cubegroup/devops-java

#image devops-java-node
cd $path/shell/java-node
docker build --platform=linux/amd64 -t devops-java-node .
docker tag devops-java-node cubegroup/devops-java-node
docker push cubegroup/devops-java-node