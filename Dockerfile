FROM golang:1.16-stretch as builder

COPY . /go/src/github.com/lichuan0620/secret-keeper-backend

WORKDIR /go/src/github.com/lichuan0620/secret-keeper-backend
RUN make build

FROM debian:stretch-slim

RUN mkdir -p /secret-keeper && \
    chown -R nobody:nogroup /secret-keeper

COPY --from=builder \
  /go/src/github.com/lichuan0620/secret-keeper-backend/bin/secret-keeper \
  /usr/local/bin/secret-keeper

ENV LISTEN_ADDRESS           ":8080"
ENV TELEMETRY_LISTEN_ADDRESS ":8081"

USER        nobody
WORKDIR     /secret-keeper
ENTRYPOINT  ["secret-keeper"]