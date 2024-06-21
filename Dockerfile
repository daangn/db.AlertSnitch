FROM alpine:latest

COPY . .

EXPOSE 9567

ENTRYPOINT [ "/alertsnitch" ]
