FROM --platform=linux/amd64 debian:stretch
USER root
ENV APP_PATH /go
ENV WEB_SERVER "0.0.0.0:80"

COPY ./sources.list /etc/apt/sources.list
COPY ./config.yaml $APP_PATH/local/config.yaml
COPY ./gotty_linux_amd64.tar.gz $APP_PATH/local/gotty_linux_amd64.tar.gz

RUN cd $APP_PATH && \
apt-get update -y && \
apt-get install -y curl && \
apt-get install -y tzdata sshpass && \
tar -zxvf local/gotty_linux_amd64.tar.gz >> /dev/null && \
mv gotty /usr/local/bin/gotty && \
chmod +x /usr/local/bin/gotty && \
mkdir -p ~/.ssh && \
chmod -R 600 ~/.ssh && \
echo "Asia/Shanghai" > /etc/timezone && \
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
rm -rf local/gotty_linux_amd64.tar.gz

WORKDIR $APP_PATH