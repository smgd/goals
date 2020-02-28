FROM golang:alpine
RUN apk add --no-cache git
RUN apk add --no-cache make
ADD . /go/src/goals
RUN set -ex && \
  cd /go/src/goals && \
  /usr/bin/make -f Makefile build
ENTRYPOINT ["/go/src/goals/caribou"]
CMD ["-config-path", "/go/src/goals/configs/server.toml"]
