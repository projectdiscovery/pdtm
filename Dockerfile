FROM golang:1.25-alpine AS builder
ARG VERSION
RUN apk add --no-cache git gcc musl-dev
WORKDIR /app
COPY . /app
RUN go mod download
RUN go build -ldflags "-s -w ${VERSION:+-X github.com/projectdiscovery/pdtm/internal/runner.version=$VERSION}" ./cmd/pdtm

FROM alpine:latest
RUN apk add --no-cache bind-tools ca-certificates
COPY --from=builder /app/pdtm /usr/local/bin/

ENTRYPOINT ["pdtm"]