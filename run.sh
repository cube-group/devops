#!/bin/sh -xe

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app
docker rm -f devops
docker-compose stop devops
docker-compose up -d
docker logs -f devops