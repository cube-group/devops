#!/bin/sh

echo $1
echo $2
#ssh
mkdir -p /root/.ssh
echo $1 > /root/.ssh/id_rsa
chmod 700 /root/.ssh
chmod 600 /root/.ssh/id_rsa

#docker-entrypoint.sh
echo $2 > /run.sh

#run
sh /run.sh