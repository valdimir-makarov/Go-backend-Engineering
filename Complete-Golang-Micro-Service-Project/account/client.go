package account

import (
	"context"
	"log"
)

type Account struct {
	ID   string
	Name string
}

type Client struct{}

// NewClient method now simply returns a placeholder client and does not use gRPC.
func NewClient(url string) (*Client, error) {
	// Placeholder for establishing a connection, now just returning a dummy client
	return &Client{}, nil
}

// Close method is a placeholder for closing the connection if needed
func (c *Client) Close() {
	// Placeholder logic for closing a connection, if applicable
	log.Println("Client connection closed")
}

// PostAccount method is now just a stub that logs an error message
func (c *Client) PostAccount(ctx context.Context, name string) (*Account, error) {
	log.Fatal("PostAccount method is not implemented")
	return nil, nil
}

// GetAccount method is now just a stub that logs an error message
func (c *Client) GetAccount(ctx context.Context, id string) (*Account, error) {
	log.Fatal("GetAccount method is not implemented")
	return nil, nil
}

// GetAccounts method is now just a stub that logs an error message
func (c *Client) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	log.Fatal("GetAccounts method is not implemented")
	return nil, nil
}
