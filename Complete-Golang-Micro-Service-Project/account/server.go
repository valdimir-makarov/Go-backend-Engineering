package account

import (
	"context"
	"fmt"
	"net"

	account "github.com/valdimir-makarov/Go-backend-Engineering/Complete-Golang-Micro-Service-Project/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
}
type Server struct {
	service *Service
	account.UnimplementedAccountServiceServer
}

type ServerFunc interface {
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
	ListenToGRPCServer(s Service, port int) error
}

func ListenToGRPCServer(s *Service, port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	serv := grpc.NewServer()

	account.RegisterAccountServiceServer(serv, &Server{service: s})
	reflection.Register(serv)
	return serv.Serve(lis)

}

func (s *Server) PostAccount(ctx context.Context, request *account.PostAccountRequest) (*account.PostAccountResponse, error) {

	accountServiceVar, err := s.service.PostAccount(ctx, request.Name)
	if err != nil {
		return nil, err
	}
	return &account.PostAccountResponse{
		Account: &account.Account{

			Name: accountServiceVar.Name,
		},
	}, nil

}
func (s *Server) GetAccounts(ctx context.Context, request *account.GetAccountsRequest) (*account.GetAccountsResponse, error) {

	res, err := s.service.GetAccounts(ctx, request.Skip, request.Take)
	if err != nil {
		return nil, err
	}
	accounts := []*account.Account{}
	for _, p := range res {
		accounts = append(
			accounts,
			&account.Account{
				Id:   p.ID,
				Name: p.Name,
			},
		)
	}
	return &account.GetAccountsResponse{Accounts: accounts}, nil

}
