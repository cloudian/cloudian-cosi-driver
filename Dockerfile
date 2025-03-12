FROM scratch

COPY build/cosi-driver /cosi-driver

ENTRYPOINT [ "/cosi-driver" ]