package account

import (
	"context"
	"log"

	account "github.com/valdimir-makarov/Go-backend-Engineering/Complete-Golang-Micro-Service-Project/account/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	connection *grpc.ClientConn
	Service    account.AccountServiceClient
}

// NewClient method now simply returns a placeholder client and does not use gRPC.
func NewClient(url string) (*Client, error) {

	connection, err := grpc.NewClient(url, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("failed TO connect to Account GRPC Client")
	}
	client := account.NewAccountServiceClient(connection)

	return &Client{
		Service:    client,
		connection: connection,
	}, nil

}

// Close method is a placeholder for closing the connection if needed
func (c *Client) Close() {

	c.connection.Close()
	// Placeholder logic for closing a connection, if applicable
	log.Println("Client connection closed")
}

// PostAccount method is now just a stub that logs an error message
func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	r, err := c.Service.PostAccount(
		ctx,
		&account.PostAccountRequest{Name: name},
	)
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:   r.Account.Id,
		Name: r.Account.Name,
	}, nil
}

func (c *Client) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	r, err := c.Service.GetAccounts(
		ctx,
		&account.GetAccountsRequest{
			Skip: skip,
			Take: take,
		},
	)
	if err != nil {
		return nil, err
	}
	accounts := []Account{}
	for _, a := range r.Accounts {
		accounts = append(accounts, Account{
			ID:   a.Id,
			Name: a.Name,
		})
	}
	return accounts, nil
}
