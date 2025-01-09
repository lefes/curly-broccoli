package repository

import (
	"database/sql"

	"github.com/lefes/curly-broccoli/internal/domain"
)

type TransactionRepo struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) *TransactionRepo {
	return &TransactionRepo{db: db}
}

func (r *TransactionRepo) CreateTransaction(transaction *domain.Transaction) error {
	// ...
	return nil
}

func (r *TransactionRepo) GetTransactionsByUserID(userID int) ([]*domain.Transaction, error) {
	// ...
	return nil, nil
}

func (r *TransactionRepo) GetAllTransactions() ([]*domain.Transaction, error) {
	// ...
	return nil, nil
}
