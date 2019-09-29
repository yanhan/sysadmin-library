FROM golang:1.12.9-buster AS builder
WORKDIR /go/src/tls-rds
COPY main.go /go/src/tls-rds/main.go
RUN go get -u github.com/go-sql-driver/mysql
RUN CGO_ENABLED=0 go build -o main

FROM alpine:3.10.2
RUN adduser -h /home/golang -s /bin/sh -D -u 2222 golang
COPY --from=builder --chown=2222:2222 /go/src/tls-rds/main /home/golang/main
COPY --chown=root:root ./rds-combined-ca-bundle.pem /etc/ssl/certs/rds-combined-ca-bundle.pem
RUN chmod 444 /etc/ssl/certs/rds-combined-ca-bundle.pem
WORKDIR /home/golang
USER 2222
ENTRYPOINT ["/home/golang/main"]
