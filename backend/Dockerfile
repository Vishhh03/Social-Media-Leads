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

# Copy binary from builder
COPY --from=builder /build/server .

# Copy migrations
COPY --from=builder /build/internal/store/migrations ./internal/store/migrations

EXPOSE 8080

CMD ["./server"]
