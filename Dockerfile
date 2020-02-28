FROM golang:alpine
RUN apk add --no-cache git
ADD . /go/src/goals
RUN set -ex && \
  cd /go/src/goals && \
  go get -d ./... && \
  go build -o caribou cmd/caribou/main.go
ENTRYPOINT ["/go/src/goals/caribou"]
CMD ["-config-path", "/go/src/goals/configs/server.toml"]
