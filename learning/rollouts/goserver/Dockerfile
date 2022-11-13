FROM golang:1.19.3-buster AS builder

WORKDIR /apps/goserver
ADD . ./
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o ./goserver -a -ldflags '-extldflags "static"'

FROM alpine:3.16.3

WORKDIR /apps/goserver
COPY --from=builder /apps/goserver/goserver ./goserver
RUN chmod a+x ./goserver
ENTRYPOINT ["./goserver"]
