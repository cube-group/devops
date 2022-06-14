FROM --platform=linux/amd64 devops-shell-java-node
ENV APP_PATH /go
ENV WEB_SERVER "0.0.0.0:80"

COPY ./web/view $APP_PATH/web/view
COPY ./web/public $APP_PATH/web/public
COPY ./bin/app $APP_PATH/app

WORKDIR $APP_PATH

ENTRYPOINT ["./app"]