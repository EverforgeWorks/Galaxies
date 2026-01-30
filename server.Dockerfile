# Step 1: Build Stage
FROM golang:1.25-alpine AS builder

# Install git for fetching dependencies
RUN apk add --no-cache git

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source
COPY . .

# Build the Main Server
RUN go build -o server ./cmd/server/main.go

# Step 2: Runtime Stage
FROM alpine:latest
RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copy Server Binary
COPY --from=builder /app/server .

# CRITICAL: Copy the static data and migrations
# Preserves the "internal/..." path structure so main.go can find them
COPY --from=builder /app/internal/data/universe.yaml ./internal/data/universe.yaml
COPY --from=builder /app/internal/adapter/repository/migrations ./internal/adapter/repository/migrations

EXPOSE 8080
CMD ["./server"]
