# Use official Golang image as a builder
FROM golang:1.23 AS builder

WORKDIR /app

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary with static linking for Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service main.go

# Use a minimal image for running
FROM alpine:3.18

WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/auth-service .

# Expose the port
EXPOSE 8080

# Command to run
CMD ["./auth-service"]
