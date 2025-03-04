package domain

import "time"

// Transaction представляет банковскую транзакцию
type Transaction struct {
    ID        int       `json:"id"`         //!< Уникальный идентификатор транзакции
    FromAccount string  `json:"from_account"` //!< Счет отправителя
    ToAccount   string  `json:"to_account"`   //!< Счет получателя
    Amount     float64   `json:"amount"`      //!< Сумма транзакции
    Currency   string    `json:"currency"`    //!< Валюта транзакции
    CreatedAt  time.Time `json:"created_at"`  //!< Время создания транзакции
}

// TransactionRepository определяет интерфейс для работы с хранилищем транзакций
type TransactionRepository interface {
    //! Создает новую транзакцию
    //! @param transaction Указатель на структуру Transaction
    //! @return error Возвращает ошибку, если операция не удалась
    Create(transaction *Transaction) error
    
    //! Возвращает транзакцию по её ID
    //! @param id Идентификатор транзакции
    //! @return *Transaction Указатель на найденную транзакцию
    //! @return error Возвращает ошибку, если транзакция не найдена
    GetByID(id int) (*Transaction, error)
    
    //! Возвращает список всех транзакций
    //! @return []*Transaction Список транзакций
    //! @return error Возвращает ошибку, если операция не удалась
    GetAll() ([]*Transaction, error)
    
    //! Обновляет существующую транзакцию
    //! @param transaction Указатель на структуру Transaction
    //! @return error Возвращает ошибку, если операция не удалась
    Update(transaction *Transaction) error
    
    //! Удаляет транзакцию по её ID
    //! @param id Идентификатор транзакции
    //! @return error Возвращает ошибку, если операция не удалась
    Delete(id int) error
} 