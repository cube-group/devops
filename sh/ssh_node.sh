#!/bin/sh

echo "sshpass -p $1 ssh -t -o 'StrictHostKeyChecking=no' -p $2 $3@$4 '$shell'"
sshpass -p $1 ssh -t -o 'StrictHostKeyChecking=no' -p $2 $3@$4 '$shell'
