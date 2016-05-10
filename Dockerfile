FROM alpine:3.3
MAINTAINER Stefano Da Ros "sd@cip.li"

WORKDIR /app
ENV HOME /app

RUN apk add --no-cache ca-certificates && \
    apk add --no-cache tzdata && \
    mkdir -p rem .config/rem
EXPOSE 42888
COPY rem.conf.example .config/rem/rem.conf
COPY rem /usr/local/bin/rem
ENTRYPOINT rem
