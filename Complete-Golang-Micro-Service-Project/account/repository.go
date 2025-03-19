package account

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Repository interface {
	Close() error
	PutAccount(ctx context.Context, account Account) error
	GetAccountById(ctx context.Context, id string) (*Account, error)
	ListAccount(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type PostgressDB struct {
	Db *sql.DB
}

// ✅ Fix: Use a pointer receiver
func (p *PostgressDB) Close() error {
	return p.Db.Close()
}

// ✅ Fix: Return a pointer instead of a value
func NewPostGressRepository(PsDbUrl string) (Repository, error) {
	db, err := sql.Open("postgres", PsDbUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	return &PostgressDB{Db: db}, nil // ✅ Return a pointer
}

func (psr *PostgressDB) PutAccount(ctx context.Context, account Account) error {
	insertQuery := `INSERT INTO account (Id, Name) VALUES ($1, $2)`
	_, err := psr.Db.ExecContext(ctx, insertQuery, account.Name)
	return err
}

func (psr *PostgressDB) GetAccountById(ctx context.Context, id string) (*Account, error) {
	query := `SELECT Id, Name FROM account WHERE Id = $1`
	var account Account
	err := psr.Db.QueryRowContext(ctx, query, id).Scan(&account.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (psr *PostgressDB) ListAccount(ctx context.Context, skip uint64, take uint64) ([]Account, error) {
	query := `SELECT Id, Name FROM account LIMIT $1 OFFSET $2`
	rows, err := psr.Db.QueryContext(ctx, query, take, skip)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allAccounts []Account

	for rows.Next() {
		var account Account
		if err := rows.Scan(&account.Name); err != nil {
			return nil, err
		}
		allAccounts = append(allAccounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return allAccounts, nil
}
