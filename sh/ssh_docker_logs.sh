#!/bin/sh

echo "sshpass -p $1 ssh -t -o 'StrictHostKeyChecking=no' -p $2 $3@$4 docker logs -f -n 1000 $5"
sshpass -p $1 ssh -t -o 'StrictHostKeyChecking=no' -p $2 $3@$4 docker logs -f -n 1000 $5
