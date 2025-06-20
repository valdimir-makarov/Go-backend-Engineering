FROM golang:1.23

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o api-gateway ./cmd/API-Gate-Way

EXPOSE 9001

CMD ["./api-gateway"]
