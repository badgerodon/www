FROM golang:1.9-alpine as builder
RUN apk --no-cache add musl-dev build-base
WORKDIR /go/src/app
COPY . .
RUN go build -o /bin/app .

FROM alpine:3.6
RUN apk --no-cache add ca-certificates
WORKDIR /root
COPY --from=0 /bin/app /root/app
COPY ./assets /root/assets
COPY ./tpl /root/tpl
CMD ["./app"]

EXPOSE 80
