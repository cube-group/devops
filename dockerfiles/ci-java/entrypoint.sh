#!/bin/sh

if [ -d "/root/.ssh" ];then
     chmod 700 /root/.ssh
fi

if [ -f "/root/.ssh/id_rsa" ];then
     chmod 600 /root/.ssh/id_rsa
fi

if [ -f "/run.sh" ];then
    sh -e /run.sh
else
    echo "no /run.sh"
fi