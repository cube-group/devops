FROM --platform=linux/amd64 golang:1.17-stretch AS builder
ENV APP_PATH /workspace
ENV GOPROXY https://goproxy.cn
COPY . $APP_PATH
RUN cd $APP_PATH && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM --platform=linux/amd64 debian:stretch
USER root
ENV APP_PATH /go
ENV WEB_SERVER "0.0.0.0:80"

COPY local/sources.list /etc/apt/sources.list
COPY ./local $APP_PATH/local
COPY ./web/view $APP_PATH/web/view
COPY ./web/public $APP_PATH/web/public
COPY --from=builder /workspace/app $APP_PATH/app

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

ENTRYPOINT ["./app"]