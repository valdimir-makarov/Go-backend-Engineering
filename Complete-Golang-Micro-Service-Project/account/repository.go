package account

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type Repository interface {
	Close()
	PutAccount(ctx context.Context, account Account) error
	GetAccountById(ctx context.Context, id string) (*Account, error)
	ListAccount(ctx context.Context, skip uint64, take uint64) ([]Account, error)
}

type PostgressDB struct {
	Db *sql.DB
}

// docker exec -ti go-micro-service  createdb -U postgres bubun

// first we are creating instance of the repository
func NewPostGressRepository(PsDbUrl string) (PostgressDB, error) {
	// Open a connection to the PostgreSQL database using the provided URL.
	db, err := sql.Open("postgres", PsDbUrl)
	if err != nil {
		return PostgressDB{}, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Verify the connection by pinging the database.
	if err := db.Ping(); err != nil {
		// Close the database connection to avoid resource leaks.
		db.Close()
		return PostgressDB{}, fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Return the initialized PostgressDB instance.
	return PostgressDB{Db: db}, nil
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

	// Check for errors from iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return allAccounts, nil
}
