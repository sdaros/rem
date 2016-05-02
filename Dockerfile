FROM alpine:3.3
MAINTAINER Stefano Da Ros "sd@cip.li"

WORKDIR /app
ENV HOME /app

RUN apk add --no-cache ca-certificates && \
    apk add --no-cache coreutils && \
    apk add --no-cache tzdata
RUN cp /usr/share/zoneinfo/Europe/Berlin /etc/localtime
RUN echo "Europe/Berlin" > /etc/timezone
RUN apk del tzdata
RUN mkdir -p rem .config/rem
COPY rem.conf .config/rem/rem.conf
COPY create.html rem/create.html
COPY rem /usr/local/bin/rem

EXPOSE 42888
CMD rem
