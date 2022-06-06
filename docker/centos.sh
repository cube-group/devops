#!/bin/sh
yum remove docker \
                  docker-client \
                  docker-client-latest \
                  docker-common \
                  docker-latest \
                  docker-latest-logrotate \
                  docker-logrotate \
                  docker-engine
# 1.  安装Docker的依赖库。
yum install -y yum-utils device-mapper-persistent-data lvm2
# 2.  添加Docker CE的软件源信息。
yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
# 3.安装Docker CE。
yum makecache fast
yum -y install docker-ce
# 4.  启动Docker服务。
systemctl start docker
systemctl enable docker