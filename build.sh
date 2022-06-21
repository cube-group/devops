#!/bin/sh -e

path=$(pwd)

#build ci
docker build --platform=linux/amd64 -t devops-ssh -f $path/dockerfiles/ssh/Dockerfile $path/dockerfiles/ssh
docker build --platform=linux/amd64 -t devops-ci-java -f $path/dockerfiles/ci-java/Dockerfile $path/dockerfiles/ci-java
#build shell
docker build --platform=linux/amd64 -t devops-shell -f $path/dockerfiles/shell/Dockerfile $path/dockerfiles

#upload
docker tag devops-ci cubegroup/devops-ssh
docker push cubegroup/devops-ssh
docker tag devops-ci-java cubegroup/devops-ci-java
docker push cubegroup/devops-ci-java

##build app
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/app
#
##build devops
#docker build --platform=linux/amd64 -t devops $path
#docker tag devops cubegroup/devops
#docker push cubegroup/devops
#
#docker build --platform=linux/amd64 -t devops-java -f $path/dockerfiles/devops-java/Dockerfile $path
#docker tag devops-java cubegroup/devops-java
#docker push cubegroup/devops-java