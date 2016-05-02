FROM alpine:3.3
MAINTAINER Stefano Da Ros "sd@cip.li"

WORKDIR /app
ENV HOME /app

RUN apk add --no-cache ca-certificates && \
    apk add --no-cache coreutils && \
    apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Europe/Berlin /etc/localtime && \
    echo "Europe/Berlin" > /etc/timezone && \
    apk del tzdata && \
    mkdir -p rem .config/rem
EXPOSE 42888
COPY rem.conf .config/rem/rem.conf
COPY create.html rem/create.html
COPY rem /usr/local/bin/rem
ENTRYPOINT rem
