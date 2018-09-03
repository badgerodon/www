FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root
COPY ./bin/badgerodon-www /root/badgerodon-www
COPY ./assets /root/assets
COPY ./tpl /root/tpl
CMD ["./badgerond-www"]

EXPOSE 80
