# syntax=docker/dockerfile:1
# Build stage
# NOTE: コンテキストのルートがrepository rootになるようにdocker buildコマンドを実行することを想定している
# go.modやgo.sum、pkgなど共通ファイルを正しくマウントするため

FROM golang:1.25-alpine AS builder

# Set working directory
WORKDIR /app

# Download dependencies
RUN --mount=type=cache,target=/go/pkg/mod/,sharing=locked \
    --mount=type=bind,source=go.sum,target=/app/go.sum \
    --mount=type=bind,source=go.mod,target=/app/go.mod \
    go mod download -x

# Build the command service binary
# Copy the entire workspace (needed for internal package references)
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=.,target=/app/ \
    CGO_ENABLED=0 GOOS=linux \
    go build -a -installsuffix cgo -o /bin/main /app/service/command/cmd/server

# Runtime stage
FROM alpine:latest AS runtime

# Install CA certificates, tzdata, and grpc-health-probe
RUN --mount=type=cache,target=/var/cache/apk,sharing=locked \
    apk --no-cache add ca-certificates tzdata curl wget && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.41/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

# Set timezone
ENV TZ=Asia/Tokyo

# Create non-root user
RUN addgroup -g 1000 app && \
    adduser -D -u 1000 -G app app

WORKDIR /app

# Copy binary from builder
COPY --from=builder /bin/main .
# Copy config from build context
COPY ./service/command/config.toml .

# Change ownership
RUN chown -R app:app /app

# Switch to non-root user
USER app

# Expose port (default command service port)
EXPOSE 80

# Run the binary
CMD ["./main"]
