FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Build the Main Server
RUN go build -o server ./cmd/server/main.go

# Uncomment the following line and in COPY if you need to reseed the universe with fresh systems.
RUN go build -o seeder ./cmd/seeder/main.go

# Step 2: Runtime container
FROM alpine:latest
WORKDIR /root/

# Copy the Server
COPY --from=builder /app/server .
COPY --from=builder /app/seeder .


EXPOSE 8080
CMD ["./server"]