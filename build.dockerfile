FROM alpine:3.5
RUN apk --no-cache add build-base go musl-dev git xz

ENV GOPATH /root
WORKDIR /root/src/github.com/badgerodon/www
COPY assets assets
COPY tpl tpl
COPY vendor vendor
COPY main.go main.go

RUN go build -o badgerodon-www .
RUN tar -cvJf /tmp/badgerodon-www.tar.xz ./tpl ./assets ./badgerodon-www
RUN ls /tmp
