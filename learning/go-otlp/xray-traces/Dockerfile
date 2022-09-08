FROM golang:1.18.6-bullseye AS builder

WORKDIR /build
COPY . ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build

FROM alpine:3.16.2

WORKDIR /app
COPY --from=builder /build/go-otlp-xray .
RUN chmod a+x ./go-otlp-xray
ENTRYPOINT ["./go-otlp-xray"]
