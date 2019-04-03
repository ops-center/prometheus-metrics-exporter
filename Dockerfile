# stage 1: build
FROM golang:1.10-alpine AS builder
LABEL maintainer="nightfury1204"

# Add source code
RUN mkdir -p /go/src/github.com/nightfury1204/prometheus-remote-metric-writer
ADD . /go/src/github.com/nightfury1204/prometheus-remote-metric-writer

# Build binary
RUN cd /go/src/github.com/nightfury1204/prometheus-remote-metric-writer && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/prometheus-remote-metric-writer

# stage 2: lightweight "release"
FROM alpine:latest
LABEL maintainer="nightfury1204"

COPY --from=builder /go/bin/prometheus-remote-metric-writer /bin/

ENTRYPOINT [ "/bin/prometheus-remote-metric-writer" ]
