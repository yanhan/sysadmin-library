FROM golang:1.13.6-buster AS builder
WORKDIR /go/src/vault-mysql
COPY main.go /go/src/vault-mysql/main.go
RUN go get -u github.com/go-sql-driver/mysql
RUN CGO_ENABLED=0 go build -o main

FROM alpine:3.11.3
RUN adduser -h /home/golang -s /bin/sh -D -u 2222 golang
COPY --from=builder --chown=2222:2222 /go/src/vault-mysql/main /home/golang/main
WORKDIR /home/golang
USER 2222
ENTRYPOINT ["/home/golang/main"]
