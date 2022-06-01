FROM --platform=linux/amd64 golang:1.17-stretch AS builder
ENV APP_PATH /workspace
ENV GOPROXY https://goproxy.cn
COPY . $APP_PATH
RUN cd $APP_PATH && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/app

FROM --platform=linux/amd64 devops-shell
USER root
ENV APP_PATH /go
ENV WEB_SERVER "0.0.0.0:80"

COPY ./config.yaml $APP_PATH/config.yaml
COPY ./web/view $APP_PATH/web/view
COPY ./web/public $APP_PATH/web/public
COPY --from=builder /workspace/bin/app $APP_PATH/app

WORKDIR $APP_PATH

ENTRYPOINT ["./app"]