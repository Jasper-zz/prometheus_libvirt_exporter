FROM       alpine:3.14

MAINTAINER Jasper

RUN        sed  -i 's#https://dl-cdn.alpinelinux.org#http://mirrors.ustc.edu.cn/#' /etc/apk/repositories && \
           apk add --no-cache libvirt-dev && \
           mkdir /lib64 && \
           ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY       ./.build/prometheus_libvirt_exporter /usr/bin

EXPOSE     9108
ENTRYPOINT ["/usr/bin/prometheus_libvirt_exporter"]
