# chat-service/chat-service/Dockerfile

FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o chat-service ./cmd

# Final stage
FROM debian:bullseye-slim

WORKDIR /root/

COPY --from=builder /app/chat-service .
COPY --from=builder /app/.env .

EXPOSE 3001

CMD ["./chat-service"]
