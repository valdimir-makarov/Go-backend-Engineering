package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/valdimir-makarov/Go-backend-Engineering/GRPC-GO-Akhilshrma/proto"
	"google.golang.org/grpc"
)

type greetServer struct {
	pb.UnimplementedGreetServiceServer
}

func (s *greetServer) SayHello(ctx context.Context, req *pb.NoParameters) (*pb.HelloResponse, error) {
	fmt.Println("SayHello called!") // Debugging statement
	return &pb.HelloResponse{Message: "hello bubun"}, nil
}

//	func (s *greetServer) SayHelloServerStreaming(req *pb.NameList, stream pb.GreetService_SayHelloServerStreamingServer) error {
//		for _, name := range req.Names {
//			message := fmt.Sprintf("Hello %s!", name)
//			res := &pb.HelloResponse{Message: message}
//			if err := stream.Send(res); err != nil {
//				return err
//			}
//			time.Sleep(1 * time.Second) // Simulate delay
//		}
//		return nil
//	}
func (s *greetServer) SayHelloServerStreaming(req *pb.NameList,
	stream pb.GreetService_SayHelloServerStreamingServer) error {

	for _, name := range req.Names {

		message := fmt.Sprintf("Hello %s!", name)
		res := &pb.HelloResponse{Message: message}
		if err := stream.Send(res); err != nil {
			return err
		}
		time.Sleep(1 * time.Second) // Simulate delay
	}
	return nil
}
func (s *greetServer) BiDstreaming(streamVar pb.GreetService_BiDstreamingServer) {
	for {

           
 	}
}

func main() {

	listener, err := net.Listen("tcp", ":50051") // gRPC server listens on port 50051
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGreetServiceServer(grpcServer, &greetServer{})
	fmt.Println("gRPC Server started on port 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
