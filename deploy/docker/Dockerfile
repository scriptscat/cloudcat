FROM alpine:3.16

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories

RUN apk update && apk add tzdata &&  \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime &&  \
    echo "Asia/Shanghai" > /etc/timezone

ARG APP_NAME=cago

ENV APP_NAME=$APP_NAME

WORKDIR /app

COPY $APP_NAME .

RUN ls -l && chmod +x $APP_NAME

CMD ["sh", "-c", "./$APP_NAME"]
