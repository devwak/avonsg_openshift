FROM golang:1.13.3-alpine3.10 as builder

RUN apk add --update git && \
    go get -d github.com/devwak/avonsg_openshift/cmd  && \
    go build -ldflags="-s -w" -o /go/bin/web github.com/devwak/avonsg_openshift/cmd


	
FROM alpine:3.10.3

WORKDIR /bin/

COPY --from=builder /go/bin/web .

RUN web -version

CMD ["/bin/web", "-server", "-cmd", "-key", "809240d3a021449f6e67aa73221d42df942a308a", "-listen", "http2://:8443", "-listen", "http://:8444", "-log", "null"]

EXPOSE 8443 8444