FROM alpine:latest

LABEL org.opencontainers.image.authors="ProjectDiscovery"
LABEL org.opencontainers.image.description="pdtm is a simple and easy-to-use golang based tool for managing open source projects from ProjectDiscovery."
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.title="pdtm"
LABEL org.opencontainers.image.url="https://github.com/projectdiscovery/pdtm"

# gcompat: purego (dlopen) makes the amd64/arm64 binaries dynamically linked
# against the glibc loader even with CGO_ENABLED=0.
RUN apk add --no-cache bind-tools ca-certificates gcompat

ARG TARGETPLATFORM
COPY $TARGETPLATFORM/pdtm /usr/local/bin/

ENTRYPOINT ["pdtm"]
