package account

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repository Repository
}

type ServiceInterface interface {
	PostAccount(ctx context.Context, name string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func InitializeService(repo Repository) *Service {
	return &Service{repository: repo}
}

func (s *Service) PostAccount(ctx context.Context, name string) (*Account, error) {
	account := &Account{
		Name: name,
		ID:   uuid.New().String(),
	}

	if err := s.repository.PutAccount(ctx, *account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *Service) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListAccount(ctx, skip, take)
}
