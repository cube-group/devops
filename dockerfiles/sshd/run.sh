#!/usr/bin/env bash

echo "root:$ROOT_PASSWORD" | chpasswd
/usr/sbin/sshd -D
