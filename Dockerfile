FROM golang:alpine
RUN apk add --no-cache git
ADD ./main.go /go/src/goals/main.go
ADD ./app /go/src/goals/app
ADD ./models /go/src/goals/models
RUN set -ex && \
  cd /go/src/goals && \
  go get -d ./... && \
  go build && \
  mv ./goals /usr/bin/goals

ENTRYPOINT [ "goals" ]
