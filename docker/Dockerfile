# Docker image for gmond2fluentd
#
# @see:      https://github.com/progrium/busybox
# @see-also: https://github.com/gliderlabs/docker-alpine/issues/12
#
#

# pull base image
FROM scratch
#FROM busybox
#FROM progrium/busybox
#FROM alpine:3.1
#FROM debian:jessie

MAINTAINER William Yeh <william.pjyeh@gmail.com>

#--
# Busybox rootfs: https://github.com/jpetazzo/docker-busybox
# ENV BUSYBOX_ROOTFS  https://github.com/jpetazzo/docker-busybox/raw/master/rootfs.tar

ADD busybox-rootfs.tar /

WORKDIR /opt
COPY    gmond2fluentd_static_linux-amd64  /opt/gmond2fluentd
COPY    wrapper                           /opt/

ENTRYPOINT ["/opt/wrapper"]
CMD ["--help"]