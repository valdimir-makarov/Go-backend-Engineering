# --- Builder Stage ---
    FROM golang:1.23-bullseye AS builder

    WORKDIR /app
    
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    
    # Verify main.go exists and build
    RUN test -f cmd/main.go || (echo "cmd/main.go not found" && exit 1)
    RUN GOOS=linux GOARCH=amd64 go build -o chat-service ./cmd
    
    # --- Final Stage ---
    FROM debian:bullseye-slim
    
    WORKDIR /app
    
    RUN apt-get update && apt-get install -y ca-certificates netcat && rm -rf /var/lib/apt/lists/*
    
    COPY --from=builder /app/chat-service .
    
    CMD ["./chat-service"]