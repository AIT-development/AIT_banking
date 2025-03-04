package repository

import (
	"bank-transactions/internal/domain"
	"database/sql"
	"time"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) domain.TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *domain.Transaction) error {
	query := `INSERT INTO transactions (from_account, to_account, amount, currency, created_at) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRow(query, 
		transaction.FromAccount,
		transaction.ToAccount,
		transaction.Amount,
		transaction.Currency,
		time.Now(),
	).Scan(&transaction.ID)
}

func (r *transactionRepository) GetByID(id int) (*domain.Transaction, error) {
	query := `SELECT id, from_account, to_account, amount, currency, created_at 
			  FROM transactions WHERE id = $1`
	var t domain.Transaction
	err := r.db.QueryRow(query, id).Scan(
		&t.ID,
		&t.FromAccount,
		&t.ToAccount,
		&t.Amount,
		&t.Currency,
		&t.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *transactionRepository) GetAll() ([]*domain.Transaction, error) {
	query := `SELECT id, from_account, to_account, amount, currency, created_at 
			  FROM transactions ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		if err := rows.Scan(
			&t.ID,
			&t.FromAccount,
			&t.ToAccount,
			&t.Amount,
			&t.Currency,
			&t.CreatedAt,
		); err != nil {
			return nil, err
		}
		transactions = append(transactions, &t)
	}
	return transactions, nil
}

func (r *transactionRepository) Update(transaction *domain.Transaction) error {
	query := `UPDATE transactions 
			  SET from_account = $1, to_account = $2, amount = $3, currency = $4 
			  WHERE id = $5`
	_, err := r.db.Exec(query,
		transaction.FromAccount,
		transaction.ToAccount,
		transaction.Amount,
		transaction.Currency,
		transaction.ID,
	)
	return err
}

func (r *transactionRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM transactions WHERE id = $1", id)
	return err
} 