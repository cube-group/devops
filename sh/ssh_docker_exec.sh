#!/bin/sh

echo "sshpass -p $1 ssh -t -o 'StrictHostKeyChecking=no' -p $2 $3@$4 docker exec -it $5 sh"
sshpass -p $1 ssh -t -o 'StrictHostKeyChecking=no' -p $2 $3@$4 docker exec -it $5 sh
