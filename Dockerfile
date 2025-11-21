# syntax=docker/dockerfile:1

FROM golang:1.24.7-alpine AS builder

ARG APP_VERSION="undefined"
ARG BUILD_TIME="undefined"

WORKDIR /go/src/github.com/artarts36/certmetrics

COPY go.mod go.sum ./
COPY pkg/collector/go.mod pkg/collector/go.sum ./pkg/collector/
COPY pkg/x509m/go.mod pkg/x509m/go.sum ./pkg/x509m/

RUN go mod download

COPY . .

RUN GOOS=linux go build -ldflags="-s -w -extldflags=-static -X 'main.Version=${APP_VERSION}' -X 'main.BuildDate=${BUILD_TIME}'" -o /go/bin/certmetrics /go/src/github.com/artarts36/certmetrics/cmd/certmetrics/main.go

######################################################

FROM alpine

WORKDIR /app

COPY --from=builder /go/bin/certmetrics /go/bin/certmetrics

# https://github.com/opencontainers/image-spec/blob/main/annotations.md
LABEL org.opencontainers.image.title="certmetrics"
LABEL org.opencontainers.image.description="prometheus exporter about certificates(x509, jwt)"
LABEL org.opencontainers.image.url="https://github.com/artarts36/certmetrics"
LABEL org.opencontainers.image.source="https://github.com/artarts36/certmetrics"
LABEL org.opencontainers.image.vendor="ArtARTs36"
LABEL org.opencontainers.image.version="$APP_VERSION"
LABEL org.opencontainers.image.created="$BUILD_TIME"
LABEL org.opencontainers.image.licenses="MIT"

EXPOSE 8080

CMD ["/go/bin/certmetrics"]
