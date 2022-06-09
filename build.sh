#!/bin/sh
#image devops-shell
cd shell
docker build --platform=linux/amd64 -t devops-shell .
docker tag devops-shell cubegroup/devops-shell
docker push cubegroup/devops-shell

#image devops
cd ..
docker build --platform=linux/amd64 -t devops .
docker tag devops cubegroup/devops
docker push cubegroup/devops

#image devops-java
cd ./shell/java
docker build --platform=linux/amd64 -t devops-java .
docker tag devops-java cubegroup/devops-java
docker push cubegroup/devops-java

#image devops-java-node
cd ../java-node
docker build --platform=linux/amd64 -t devops-java-node .
docker tag devops-java-node cubegroup/devops-java-node
docker push cubegroup/devops-java-node