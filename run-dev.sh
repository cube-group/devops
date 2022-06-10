#!/bin/sh -xe

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app
docker rm -f devops
docker-compose -f docker-compose-dev.yaml stop devops
docker-compose -f docker-compose-dev.yaml up -d
docker logs -f devops