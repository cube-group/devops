#!/bin/sh -e

path=$(pwd)
version=$(cat $path/local/version)
echo $version
docker buildx create --use --name build --node build --driver-opt network=host
build()
{
  docker buildx build --platform $1 -t $2 -f $3 $4 --push
}

build linux/arm64,linux/amd64 cubegroup/devops-shell:$version $path/dockerfiles/shell/Dockerfile $path/dockerfiles
build linux/arm64,linux/amd64 cubegroup/devops-shell-java:$version $path/dockerfiles/shell-java/Dockerfile $path/dockerfiles

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/app
build linux/arm64,linux/amd64 cubegroup/devops:$version $path/dockerfiles/devops/Dockerfile $path
build linux/arm64,linux/amd64 cubegroup/devops-java:$version $path/dockerfiles/devops-java/Dockerfile $path
#build linux/arm64,linux/amd64 cubegroup/devops:latest $path/dockerfiles/devops/Dockerfile $path
#build linux/arm64,linux/amd64 cubegroup/devops-java:latest $path/dockerfiles/devops-java/Dockerfile $path