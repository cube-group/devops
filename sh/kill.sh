#!/bin/sh

ps -A -ostat,ppid | grep -e '^[Zz]' | awk '{print $2}'