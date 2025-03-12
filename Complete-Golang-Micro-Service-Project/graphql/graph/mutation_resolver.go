package graph

import (
	"context"
	"fmt"

	"github.com/valdimir-makarov/Go-backend-Engineering/Complete-Golang-Micro-Service-Project/graphql/graph/model"
)

type mutationResolver struct{ server *Server }

func (r *mutationResolver) CreateAccount(ctx context.Context, account model.AccountInput) (*model.Account, error) {
	panic(fmt.Errorf("not implemented: CreateAccount - createAccount"))
}

// CreateProduct is the resolver for the createProduct field.
func (r *mutationResolver) CreateProduct(ctx context.Context, product model.ProductInput) (*model.Product, error) {
	panic(fmt.Errorf("not implemented: CreateProduct - createProduct"))
}

// CreateOrder is the resolver for the createOrder field.
func (r *mutationResolver) CreateOrder(ctx context.Context, order model.OrderInput) (*model.Order, error) {
	panic(fmt.Errorf("not implemented: CreateOrder - createOrder"))
}
