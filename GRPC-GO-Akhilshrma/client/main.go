package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/valdimir-makarov/Go-backend-Engineering/GRPC-GO-Akhilshrma/proto"
	"google.golang.org/grpc"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewGreetServiceClient(conn)
	req := &pb.NameList{
		Names: []string{"Bubun", "Alice", "Bob"},
	}
	stream, err := client.SayHelloServerStreaming(context.Background(), req)
	// Unary RPC Call
	if err != nil {
		log.Fatalf("Error calling SayHelloServerStreaming: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	for {
		res, err := stream.Recv()
		if err != nil {
			break // Stream finished
		}
		fmt.Println("Server Response:", res.Message)
	}
	res, err := client.SayHello(ctx, &pb.NoParameters{})
	if err != nil {
		log.Fatalf("Error calling SayHello: %v", err)
	}

	fmt.Println("Server Response:", res.Message)
}
