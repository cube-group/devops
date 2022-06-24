#!/bin/sh

if [ ! -f "/go/app" ];then
     dockerd-entrypoint.sh
else
     dockerd-entrypoint.sh &
     cd /go && ./app
fi
