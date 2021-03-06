FROM golang:alpine
RUN apk add --no-cache git make
ADD . /go/src/goals
RUN  cd /go/src/goals && /usr/bin/make -f Makefile build
ENTRYPOINT ["/go/src/goals/caribou"]
CMD ["-config-path", "/go/src/goals/configs/server.toml"]
