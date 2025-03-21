package main

import (
	"context"
	"log"
	"net"

	searchpb "your_project/search"
	productpb "your_project/products"

	"google.golang.org/grpc"
)

type searchServiceServer struct {
	searchpb.UnimplementedSearchServiceServer
	productClient productpb.ProductServiceClient
}

func (s *searchServiceServer) SearchProduct(ctx context.Context, req *searchpb.SearchRequest) (*searchpb.SearchResponse, error) {
	// Call Product Service
	productsResp, err := s.productClient.SearchProduct(ctx, &productpb.SearchRequest{Name: req.Name})
	if err != nil {
		return nil, err
	}

	// Convert response
	var result []*searchpb.Product
	for _, p := range productsResp.Products {
		result = append(result, &searchpb.Product{
			Id:    p.Id,
			Name:  p.Name,
			Price: p.Price,
		})
	}

	return &searchpb.SearchResponse{Products: result}, nil
}

func main() {
	// Connect to Product Service
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to Product Service: %v", err)
	}
	defer conn.Close()

	productClient := productpb.NewProductServiceClient(conn)

	// Start Search Service
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	searchpb.RegisterSearchServiceServer(server, &searchServiceServer{productClient: productClient})

	log.Println("Search Service running on port 50052")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
