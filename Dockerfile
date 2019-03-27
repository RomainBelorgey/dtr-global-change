FROM alpine:3.7
LABEL maintainer "RomainBelorgey"

COPY dtr-global-change /usr/bin/

ENTRYPOINT ["/usr/bin/dtr-global-change"]
