FROM golang:1.24.3-alpine AS builder
RUN apk add --no-cache git gcc musl-dev
WORKDIR /app
COPY . /app
RUN go mod download
RUN go build ./cmd/pdtm

FROM alpine:latest
RUN apk add --no-cache bind-tools ca-certificates
COPY --from=builder /app/pdtm /usr/local/bin/

ENTRYPOINT ["pdtm"]