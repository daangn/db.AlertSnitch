FROM golang:1.14 AS builder

WORKDIR /Users/silverlee/.git/db.AlertSnitch

COPY . .

RUN env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o /alertsnitch .

FROM amazonlinux:latest

COPY --from=builder /alertsnitch /alertsnitch
USER nobody
ENTRYPOINT ["/alertsnitch"]
