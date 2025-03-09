# Stage 1: Build the Go application
FROM golang:1.22 AS builder

# Set the working directory inside the container
WORKDIR /app/graphql

# Copy the Go module files from the root of the project
COPY ../go.mod ../go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
# Assuming the entry point is server.go
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./server.go
# Stage 2: Create a minimal runtime image
FROM alpine:latest

# Install necessary libraries (if any)
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/graphql/server .

# Expose the port your server listens on (e.g., 8080)
EXPOSE 8080

# Command to run the server
CMD ["./server"]