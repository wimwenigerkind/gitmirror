FROM alpine:3.23

RUN apk add --no-cache git ca-certificates

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/gitmirror /usr/local/bin/gitmirror

WORKDIR /work
ENTRYPOINT ["/usr/local/bin/gitmirror"]