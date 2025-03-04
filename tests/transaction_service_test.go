package service

import (
    "bank-transactions/internal/domain"
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) Create(transaction *domain.Transaction) error {
    args := m.Called(transaction)
    return args.Error(0)
}

func (m *MockRepository) GetByID(id int) (*domain.Transaction, error) {
    args := m.Called(id)
    return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockRepository) GetAll() ([]*domain.Transaction, error) {
    args := m.Called()
    return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockRepository) Update(transaction *domain.Transaction) error {
    args := m.Called(transaction)
    return args.Error(0)
}

func (m *MockRepository) Delete(id int) error {
    args := m.Called(id)
    return args.Error(0)
}

func TestCreateTransaction(t *testing.T) {
    mockRepo := new(MockRepository)
    service := NewTransactionService(mockRepo)

    validTransaction := &domain.Transaction{
        FromAccount: "acc1",
        ToAccount:   "acc2",
        Amount:      100.0,
    }

    // Test valid transaction
    mockRepo.On("Create", validTransaction).Return(nil)
    err := service.CreateTransaction(validTransaction)
    assert.NoError(t, err)

    // Test invalid amount
    invalidTransaction := &domain.Transaction{
        FromAccount: "acc1",
        ToAccount:   "acc2",
        Amount:      -100.0,
    }
    err = service.CreateTransaction(invalidTransaction)
    assert.Error(t, err)
    assert.Equal(t, "amount must be positive", err.Error())

    mockRepo.AssertExpectations(t)
}

