# Step 1: Build stage (Linux)
FROM golang:1.24-bullseye AS builder

WORKDIR /app

# Copy go.mod/go.sum and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy all source files
COPY . .

# Build for Linux ARM64
RUN GOOS=linux GOARCH=arm64 go build -o notification-service .

# Step 2: Runtime stage
FROM debian:bullseye-slim

WORKDIR /app

# Copy the built binary
COPY --from=builder /app/notification-service .

# Expose app port
EXPOSE 8080

# Run the binary
CMD ["./notification-service"]