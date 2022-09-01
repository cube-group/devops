#!/bin/sh -e

path=$(pwd)
version=latest
echo $version
docker buildx create --use --name build --node build --driver-opt network=host
build()
{
  docker buildx build --platform $1 -t $2 -f $3 $4 --push
}

build linux/arm64,linux/amd64 cubegroup/sshd-docker:$version $path/Dockerfile $path