FROM golang:1.13.2-buster AS builder
WORKDIR /go/src/tls-redis
COPY main.go /go/src/tls-redis/main.go
RUN go mod init github.com/yanhan/sysadmin-library/tls-redis/eks
RUN CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -o main

FROM alpine:3.10.2
RUN adduser -h /home/golang -s /bin/sh -D -u 2222 golang
COPY --from=builder --chown=2222:2222 /go/src/tls-redis/main /home/golang/main
WORKDIR /home/golang
USER 2222
ENTRYPOINT ["/home/golang/main"]
