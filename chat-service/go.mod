module github.com/valdimir-makarov/Go-backend-Engineering/chat-service

go 1.21



// version: "3.8"

// services:
//   postgres:
//     image: postgres:latest
//     container_name: my_postgres
//     environment:
//       POSTGRES_USER: myuser
//       POSTGRES_PASSWORD: mypassword
//       POSTGRES_DB: ChatAppDb
//     ports:
//       - "9002:5432"
//     volumes:
//       - pg_data:/var/lib/postgresql/data
//       - ./init-db.sql:/docker-entrypoint-initdb.d/init-db.sql
//     networks:
//       - my_network
//     healthcheck:
//       test: ["CMD-SHELL", "pg_isready -U myuser -d ChatAppDb"]
//       interval: 10s
//       timeout: 5s
//       retries: 5
//       start_period: 60s

//   pgadmin:
//     image: dpage/pgadmin4
//     container_name: my_pgadmin
//     environment:
//       PGADMIN_DEFAULT_EMAIL: admin@example.com
//       PGADMIN_DEFAULT_PASSWORD: admin
//     ports:
//       - "9003:80"
//     depends_on:
//       postgres:
//         condition: service_healthy
//     networks:
//       - my_network

//   chat-service:
//     build:
//       context: .
//       dockerfile: Dockerfile
//     container_name: chat_service
//     depends_on:
//       postgres:
//         condition: service_healthy
//     environment:
//       CHAT_SERVICE_PORT: 3001
//       DB_USER: myuser
//       DB_PASSWORD: mypassword
//       DB_NAME: ChatAppDb
//       DB_HOST: postgres
//       DB_PORT: "5432"
//     ports:
//       - "9006:3001"
//     networks:
//       - my_network
//     healthcheck:
//       test: ["CMD-SHELL", "nc -z localhost 3001 || exit 1"]
//       interval: 10s
//       timeout: 5s
//       retries: 5
//       start_period: 120s

// volumes:
//   pg_data:

// networks:
//   my_network:
//     name: my_network
//     driver: bridge



// # --- Builder Stage ---
//     FROM golang:1.23-bullseye AS builder

//     WORKDIR /app
    
//     COPY go.mod go.sum ./
//     RUN go mod download
    
//     COPY . .
    
//     # Verify main.go exists and build
//     RUN test -f cmd/main.go || (echo "cmd/main.go not found" && exit 1)
//     RUN GOOS=linux GOARCH=amd64 go build -o chat-service ./cmd
    
//     # --- Final Stage ---
//     FROM debian:bullseye-slim
    
//     WORKDIR /app
    
//     RUN apt-get update && apt-get install -y ca-certificates netcat && rm -rf /var/lib/apt/lists/*
    
//     COPY --from=builder /app/chat-service .
    
//     CMD ["./chat-service"]