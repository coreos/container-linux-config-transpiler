FROM golang:alpine
ENV CGO_ENABLED=0
RUN apk update && apk add --virtual .build-deps bash git make \
    && git clone https://github.com/coreos/container-linux-config-transpiler \
    && cd container-linux-config-transpiler \
    && make \
    && mv bin/ct /usr/bin/ \
    && rm -rf /var/cache/apk \
    && rm -rf /go/container-linux-config-transpiler
ENTRYPOINT ["/usr/bin/ct"]
CMD ["-pretty"]
