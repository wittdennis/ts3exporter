FROM golang:1.20 AS build

ADD . /go/src/github.com/hikhvar/ts3exporter

RUN cd /go/src/github.com/hikhvar/ts3exporter && \
    go get -d -v ./... && \
    CGO_ENABLED=0 go build -o /go/bin/ts3exporter

RUN mkdir -p /rootfs/etc && \
    cp /go/bin/ts3exporter /rootfs/ && \
    echo "nogroup:*:100:nobody" > /rootfs/etc/group && \
    echo "nobody:*:100:100:::" > /rootfs/etc/passwd


FROM alpine:3.18.0
COPY --from=build --chown=100:100 /rootfs /
USER 100:100
EXPOSE 9189/tcp
ENTRYPOINT ["/ts3exporter"]
