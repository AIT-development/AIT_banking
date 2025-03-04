package repository

import (
    "bank-transactions/internal/domain"
    "database/sql"
    "testing"
    "time"

    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

type TransactionRepositoryTestSuite struct {
    suite.Suite
    db   *sql.DB
    mock sqlmock.Sqlmock
    repo domain.TransactionRepository
}

func (s *TransactionRepositoryTestSuite) SetupTest() {
    var err error
    s.db, s.mock, err = sqlmock.New()
    assert.NoError(s.T(), err)
    s.repo = NewTransactionRepository(s.db)
}

func (s *TransactionRepositoryTestSuite) TearDownTest() {
    s.db.Close()
}

func TestTransactionRepositorySuite(t *testing.T) {
    suite.Run(t, new(TransactionRepositoryTestSuite))
}

func (s *TransactionRepositoryTestSuite) TestCreate() {
    transaction := &domain.Transaction{
        FromAccount: "acc1",
        ToAccount:   "acc2",
        Amount:      100.0,
        Currency:    "USD",
    }

    s.mock.ExpectQuery("INSERT INTO transactions").
        WithArgs(transaction.FromAccount, transaction.ToAccount, transaction.Amount, transaction.Currency, sqlmock.AnyArg()).
        WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

    err := s.repo.Create(transaction)
    assert.NoError(s.T(), err)
    assert.Equal(s.T(), 1, transaction.ID)
}

func (s *TransactionRepositoryTestSuite) TestGetByID() {
    expected := &domain.Transaction{
        ID:          1,
        FromAccount: "acc1",
        ToAccount:   "acc2",
        Amount:      100.0,
        Currency:    "USD",
        CreatedAt:   time.Now(),
    }

    rows := sqlmock.NewRows([]string{"id", "from_account", "to_account", "amount", "currency", "created_at"}).
        AddRow(expected.ID, expected.FromAccount, expected.ToAccount, expected.Amount, expected.Currency, expected.CreatedAt)

    s.mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

    result, err := s.repo.GetByID(1)
    assert.NoError(s.T(), err)
    assert.Equal(s.T(), expected, result)
}
