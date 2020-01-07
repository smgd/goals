FROM golang:alpine
RUN apk add --no-cache git
ADD ./main.go /go/src/goals/main.go
ADD ./app /go/src/goals/app
RUN set -ex && \
  cd /go/src/goals && \
  go get -u github.com/gorilla/mux && \
  go get -u github.com/jinzhu/gorm && \
  go get -u github.com/jinzhu/copier && \
  go get -u github.com/gorilla/handlers && \
  go get -u github.com/dgrijalva/jwt-go && \
  go get -u github.com/lib/pq && \
  go build && \
  mv ./goals /usr/bin/goals

ENTRYPOINT [ "goals" ]
