ARG ALPINE_VERSION=3.17
ARG GOLANG_VERSION=1.19.2

FROM golang:${GOLANG_VERSION} AS build-env
WORKDIR /go/src/github.com/scorum/cosmos-network/
COPY . .
RUN make linux

FROM alpine:${ALPINE_VERSION}
RUN apk update && apk add --update ca-certificates libc6-compat
COPY --from=build-env /go/src/github.com/scorum/cosmos-network/build/scorumd-linux-amd64 /usr/bin/scorumd
COPY --from=build-env /go/src/github.com/scorum/cosmos-network/genesis/testnet/ /etc/testnet/

EXPOSE 26657
EXPOSE 26656
EXPOSE 9090
EXPOSE 1317

CMD ["scorumd", "start"]