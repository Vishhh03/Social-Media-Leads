# ============================
# Build Stage
# ============================
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/server ./cmd/api

# ============================
# Production Stage
# ============================
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Create a non-root user
RUN adduser -D -g '' appuser

# Copy binary from builder
COPY --from=builder /build/server .

# Copy migrations
COPY --from=builder /build/internal/store/migrations ./internal/store/migrations

# Set ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]
