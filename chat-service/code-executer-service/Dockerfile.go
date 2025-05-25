FROM golang:1.22
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o api-server ./api/main.go
CMD ["./api-server"]