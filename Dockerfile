FROM golang:1.13-rc-alpine

RUN apk add --update \
    build-base \
    ca-certificates \
    musl-dev \
    git \
  && rm -rf /var/cache/apk/*

RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

ENV GO111MODULE=on

WORKDIR /go/src/github.com/badgerodon/www
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN go install -ldflags='-s -w' -tags netgo -installsuffix netgo -v ./...

CMD ["/go/bin/www"]
