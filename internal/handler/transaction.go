package handler

import (
    "bank-transactions/internal/domain"
    "bank-transactions/internal/service"
    "net/http"
    "strconv"
    
    "github.com/labstack/echo/v4"
)

type TransactionHandler struct {
    service *service.TransactionService
}

func NewTransactionHandler(s *service.TransactionService) *TransactionHandler {
    return &TransactionHandler{service: s}
}

func (h *TransactionHandler) Create(c echo.Context) error {
    var transaction domain.Transaction
    if err := c.Bind(&transaction); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
    }

    if err := h.service.CreateTransaction(&transaction); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusCreated, transaction)
}

func (h *TransactionHandler) GetByID(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
    }

    transaction, err := h.service.GetTransaction(id)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }
    if transaction == nil {
        return c.JSON(http.StatusNotFound, map[string]string{"error": "transaction not found"})
    }

    return c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) GetAll(c echo.Context) error {
    transactions, err := h.service.GetAllTransactions()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }
    return c.JSON(http.StatusOK, transactions)
}

func (h *TransactionHandler) Update(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
    }

    var transaction domain.Transaction
    if err := c.Bind(&transaction); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
    }
    transaction.ID = id

    if err := h.service.UpdateTransaction(&transaction); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    return c.JSON(http.StatusOK, transaction)
}

func (h *TransactionHandler) Delete(c echo.Context) error {
    id, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ID"})
    }

    if err := h.service.DeleteTransaction(id); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
    }

    return c.NoContent(http.StatusNoContent)
} 