FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata

ENV SERVER_ADDR ":4000"
EXPOSE 4000

RUN mkdir -p /usr/local/app
WORKDIR /usr/local/app

COPY ./bin/bloefish ./bloefish
ENTRYPOINT ["./bloefish"]
