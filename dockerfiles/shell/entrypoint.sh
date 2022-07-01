#!/bin/sh

if [ ! -f "/go/app" ];then
     dockerd-entrypoint.sh
else
     if [ $DOCKER_IN_DOCKER ];then
        dockerd-entrypoint.sh &
     fi
     cd /go && ./app
fi