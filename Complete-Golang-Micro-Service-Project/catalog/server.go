package catalog

import (
	"context"
	"fmt"
	"log"
	"net"

	catalog "github.com/valdimir-makarov/Go-backend-Engineering/Complete-Golang-Micro-Service-Project/catalog/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcServer struct {
	catalog.UnimplementedCatalogServiceServer

	service Service
}

// mustEmbedUnimplementedCatalogServiceServer implements catalog.CatalogServiceServer.

func ListenGRPC(s Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()
	catalog.RegisterCatalogServiceServer(serv, &grpcServer{service: s})

	reflection.Register(serv)
	return serv.Serve(lis)
}

func (s *grpcServer) PostProduct(ctx context.Context, r *catalog.PostProductRequest) (*catalog.PostProductResponse, error) {
	p, err := s.service.PostProduct(ctx, r.Name, r.Description, r.Price)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &catalog.PostProductResponse{Product: &catalog.Product{
		Id:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
	}}, nil
}

func (s *grpcServer) GetProduct(ctx context.Context, r *catalog.GetProductRequest) (*catalog.GetProductResponse, error) {
	p, err := s.service.GetProduct(ctx, r.Id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &catalog.GetProductResponse{
		Product: &catalog.Product{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
		},
	}, nil
}

func (s *grpcServer) GetProducts(ctx context.Context, r *catalog.GetProductsRequest) (*catalog.GetProductsResponse, error) {
	var res []Product
	var err error
	if r.Query != "" {
		res, err = s.service.SearchProducts(ctx, r.Query, r.Skip, r.Take)
	} else if len(r.Ids) != 0 {
		res, err = s.service.GetProductsByIDs(ctx, r.Ids)
	} else {
		res, err = s.service.GetProducts(ctx, r.Skip, r.Take)
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}

	products := []*catalog.Product{}
	for _, p := range res {
		products = append(
			products,
			&catalog.Product{
				Id:          p.ID,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
			},
		)
	}
	return &catalog.GetProductsResponse{Products: products}, nil
}
