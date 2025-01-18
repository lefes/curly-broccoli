package services

import (
	"github.com/lefes/curly-broccoli/internal/domain"
	"github.com/lefes/curly-broccoli/internal/repository"
)

type TransactionService struct {
	repo repository.Transactions
}

func NewTransactionService(repo repository.Transactions) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) CreateTransaction(transaction *domain.Transaction) error {
	return s.repo.CreateTransaction(transaction)
}

func (s *TransactionService) GetTransactionsByUserID(userID int) ([]*domain.Transaction, error) {
	return s.repo.GetTransactionsByUserID(userID)
}

func (s *TransactionService) GetAllTransactions() ([]*domain.Transaction, error) {
	return s.repo.GetAllTransactions()
}
