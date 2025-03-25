# Use the correct Go version
FROM golang:1.21.3 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project
COPY . .

# Statically compile the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o chat-service ./cmd/main.go

# Use a minimal base image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/chat-service .

# Run the application
CMD ["./chat-service"]