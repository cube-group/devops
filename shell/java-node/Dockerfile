FROM devops
VOLUME /data
#nodejs
RUN apt install -y nodejs npm && \
    mkdir -p /data/nodejs/node_modules/node_global && \
    mkdir -p /data/nodejs/node_modules/node_cache && \
    npm config set unsafe-perm true && \
    npm config set registry https://registry.npm.taobao.org && \
    npm install -g npm && \
    npm config set prefix "/data/nodejs/node_modules/node_global" && \
    npm config set cache "/data/nodejs/node_modules/node_cache" && \
    #npm-cli-adduser
    npm install npm-cli-adduser -g && \
    ln -s /data/nodejs/node_modules/node_global/lib/node_modules/npm-cli-adduser/index.js /usr/local/bin/npm-cli-adduser && \
    #cnpm
    npm install -g cnpm --registry=https://registry.npmmirror.com
