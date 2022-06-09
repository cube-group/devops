#!/bin/sh

cd shell
docker build -t devops-shell .

cd ..
docker build -t devops .