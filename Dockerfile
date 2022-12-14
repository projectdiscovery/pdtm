FROM golang:1.19.4-alpine AS builder
RUN apk add --no-cache git
RUN go install -v github.com/projectdiscovery/pdtm/cmd/pdtm@latest

FROM alpine:3.17.0
RUN apk -U upgrade --no-cache \
    && apk add --no-cache bind-tools ca-certificates
COPY --from=builder /go/bin/pdtm /usr/local/bin/

ENTRYPOINT ["pdtm"]
