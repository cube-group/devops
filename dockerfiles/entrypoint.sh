#!/bin/sh

if [ -f "/go/app" ];then
     cd /go && ./app
else
     echo "no /go/app"
     bash
fi