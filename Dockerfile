FROM golang:1.14 AS builder

WORKDIR /app

COPY . .

RUN env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o /tmp/alertsnitch .

FROM amazonlinux:latest

COPY --from=builder /tmp/alertsnitch /usr/local/bin/alertsnitch

RUN chmod +x /usr/local/bin/alertsnitch

USER nobody

ENTRYPOINT ["/usr/local/bin/alertsnitch"]
