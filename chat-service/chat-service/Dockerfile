# --- Builder Stage ---
    FROM golang:1.23-bullseye AS builder

    WORKDIR /app
    
    # Install air for live reloading
    RUN go install github.com/cosmtrek/air@v1.49.0

    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    
    # Verify main.go exists and build (optional for dev, but good check)
    RUN test -f cmd/main.go || (echo "cmd/main.go not found" && exit 1)
    # RUN GOOS=linux GOARCH=amd64 go build -o chat-service ./cmd/main.go
    # RUN ls -lh /app/chat-service

    # --- Final Stage ---
    # For development with air, we can stay in the golang image or copy air to a smaller image
    # But staying in golang image is easier for dev tools
    
    CMD ["air", "-c", ".air.toml"]