#!/bin/sh -e

path=$(pwd)
docker buildx create --use --name build --node build --driver-opt network=host
build()
{
  docker buildx build --platform $1 -t $2 -f $3 $4 --push
}

build linux/arm64,linux/amd64 cubegroup/devops-shell:v2 $path/dockerfiles/shell/Dockerfile $path/dockerfiles/shell
build linux/arm64,linux/amd64 cubegroup/devops-ssh $path/dockerfiles/ssh/Dockerfile $path/dockerfiles/ssh
build linux/amd64 cubegroup/devops-ci-java $path/dockerfiles/ci-java/Dockerfile $path/dockerfiles/ci-java


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