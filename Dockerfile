FROM golang:1.10.3-alpine
LABEL maintainer="Alex Romanin alexandr@endpass.com"

ADD . /go/src/github.com/machinae/wildproxy

WORKDIR /go/src/github.com/machinae/wildproxy

RUN apk update && \
    apk add git curl && \
    rm -rf /var/cache/apk/* && \
    go get -d -v && \
    go install -v && \
    rm -rf /go/src/*

HEALTHCHECK --interval=10s --timeout=1m --retries=5 CMD curl http://localhost:8080/health || exit 1

ENTRYPOINT ["/go/bin/wildproxy"]

EXPOSE 8080
