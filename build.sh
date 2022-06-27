#!/bin/sh -e

path=$(pwd)
docker buildx create --use --name build --node build --driver-opt network=host
build()
{
  docker buildx build --platform $1 -t $2 -f $3 $4 --push
}

build linux/arm64,linux/amd64 cubegroup/devops-shell:v2 $path/dockerfiles/shell/Dockerfile $path/dockerfiles/shell
build linux/arm64,linux/amd64 cubegroup/devops-ci-default $path/dockerfiles/ci-default/Dockerfile $path/dockerfiles/ci-default
build linux/amd64 cubegroup/devops-ci-java $path/dockerfiles/ci-java/Dockerfile $path/dockerfiles/ci-java

#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/app
#build linux/arm64,linux/amd64 cubegroup/devops:v2 $path/Dockerfile $path