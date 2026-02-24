FROM alpine:3.23.3

# TOMTODO Debug
#RUN apk add --upgrade --no-cache bash

COPY build/cosi-driver /cosi-driver

COPY docker/entrypoint.sh /entrypoint.sh

RUN apk add --upgrade --no-cache ca-certificates

ENTRYPOINT [ "/entrypoint.sh" ]
#ENTRYPOINT [ "/cosi-driver" ]
