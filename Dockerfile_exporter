# syntax=docker/dockerfile:1

FROM golang:1.23.3-alpine AS builder

ARG APP_VERSION="undefined"
ARG BUILD_TIME="undefined"

WORKDIR /go/src/github.com/artarts36/certmetrics

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux go build -ldflags="-s -w -extldflags=-static -X 'main.Version=${APP_VERSION}' -X 'main.BuildDate=${BUILD_TIME}'" -o /go/bin/certmetrics /go/src/github.com/artarts36/certmetrics/exporter/main.go

######################################################

FROM alpine

WORKDIR /app

COPY --from=builder /go/bin/certmetrics /go/bin/certmetrics

# https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.title="certmetrics-exporter"
LABEL org.opencontainers.image.description="prometheus exporter about certificates(x509, jwt)"
LABEL org.opencontainers.image.url="https://github.com/artarts36/certmetrics"
LABEL org.opencontainers.image.source="https://github.com/artarts36/certmetrics"
LABEL org.opencontainers.image.vendor="ArtARTs36"
LABEL org.opencontainers.image.version="$APP_VERSION"
LABEL org.opencontainers.image.created="$BUILD_TIME"
LABEL org.opencontainers.image.licenses="MIT"

EXPOSE 8080

CMD ["/go/bin/certmetrics"]
