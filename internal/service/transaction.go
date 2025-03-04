package service

import (
    "bank-transactions/internal/domain"
    "errors"
)

type transactionService struct {
    repo domain.TransactionRepository
}

func NewTransactionService(repo domain.TransactionRepository) *transactionService {
    return &transactionService{repo: repo}
}

func (s *transactionService) CreateTransaction(transaction *domain.Transaction) error {
    if transaction.FromAccount == "" || transaction.ToAccount == "" {
        return errors.New("from_account and to_account are required")
    }
    if transaction.Amount <= 0 {
        return errors.New("amount must be positive")
    }
    return s.repo.Create(transaction)
}

func (s *transactionService) GetTransaction(id int) (*domain.Transaction, error) {
    return s.repo.GetByID(id)
}

func (s *transactionService) GetAllTransactions() ([]*domain.Transaction, error) {
    return s.repo.GetAll()
}

func (s *transactionService) UpdateTransaction(transaction *domain.Transaction) error {
    if transaction.ID == 0 {
        return errors.New("transaction ID is required")
    }
    return s.repo.Update(transaction)
}

func (s *transactionService) DeleteTransaction(id int) error {
    return s.repo.Delete(id)
} 