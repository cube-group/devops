#!/bin/sh

if [ -f "/root/.ssh/id_rsa" ];then
     chmod 700 /root/.ssh
     chmod 600 /root/.ssh/id_rsa
fi

if [ -f "/root/.ssh2/id_rsa" ];then
     chmod 700 /root/.ssh2
     chmod 600 /root/.ssh2/id_rsa
fi

if [ -f "/.devops/run.sh" ];then
    cd /.devops && sh -e run.sh
else
    echo "no /.devops/run.sh"
fi