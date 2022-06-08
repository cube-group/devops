#!/bin/sh

ps aux | grep -e "MD5=$1" | grep -v grep | awk '{print $2}' | sort -rn | sed -n '1p' | xargs kill
