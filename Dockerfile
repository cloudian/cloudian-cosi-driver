FROM alpine:3.23.3

COPY build/cosi-driver /cosi-driver

COPY docker/entrypoint.sh /entrypoint.sh

RUN apk add --upgrade --no-cache ca-certificates

ENTRYPOINT [ "/entrypoint.sh" ]
