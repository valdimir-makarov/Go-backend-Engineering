# chat-service/chat-service/Dockerfile

FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/cosmtrek/air@latest

RUN GOOS=linux GOARCH=amd64 go build -o chat-service ./cmd
RUN go install github.com/cosmtrek/air@latest

# Final stage
FROM debian:bullseye-slim

WORKDIR /root/

COPY --from=builder /app/chat-service .
COPY --from=builder /app/.env .

EXPOSE 3001


CMD ["air", "-c", ".air.toml"]
