# Docker image for gmond2fluentd
#
# @see:      https://github.com/progrium/busybox
# @see-also: https://github.com/gliderlabs/docker-alpine/issues/12
#
#

# pull base image
FROM progrium/busybox
#FROM alpine:3.1
#FROM debian:jessie

MAINTAINER William Yeh <william.pjyeh@gmail.com>

WORKDIR /opt
COPY    gmond2fluentd_linux-amd64  /opt/gmond2fluentd
COPY    wrapper                    /opt/

ENTRYPOINT ["/opt/wrapper"]
CMD ["--help"]