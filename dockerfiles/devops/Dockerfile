FROM cubegroup/devops-shell
USER root
ENV APP_PATH /go
ENV WEB_PORT "80"

COPY ./web/view $APP_PATH/web/view
COPY ./web/public $APP_PATH/web/public
COPY ./bin/app $APP_PATH/app

VOLUME /root/.ssh
VOLUME /var/run/docker.sock