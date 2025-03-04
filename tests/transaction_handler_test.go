package handler

import (
    "bank-transactions/internal/domain"
    "bank-transactions/internal/service"
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "strconv"
    "testing"

    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockService struct {
    mock.Mock
}

func (m *MockService) CreateTransaction(transaction *domain.Transaction) error {
    args := m.Called(transaction)
    return args.Error(0)
}

func (m *MockService) GetTransaction(id int) (*domain.Transaction, error) {
    args := m.Called(id)
    return args.Get(0).(*domain.Transaction), args.Error(1)
}

func (m *MockService) GetAllTransactions() ([]*domain.Transaction, error) {
    args := m.Called()
    return args.Get(0).([]*domain.Transaction), args.Error(1)
}

func (m *MockService) UpdateTransaction(transaction *domain.Transaction) error {
    args := m.Called(transaction)
    return args.Error(0)
}

func (m *MockService) DeleteTransaction(id int) error {
    args := m.Called(id)
    return args.Error(0)
}

func TestCreateTransactionHandler(t *testing.T) {
    e := echo.New()
    mockService := new(MockService)
    handler := NewTransactionHandler(mockService)

    validTransaction := domain.Transaction{
        FromAccount: "acc1",
        ToAccount:   "acc2",
        Amount:      100.0,
    }

    // Test successful creation
    mockService.On("CreateTransaction", &validTransaction).Return(nil)
    jsonData, _ := json.Marshal(validTransaction)
    req := httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader(jsonData))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    if assert.NoError(t, handler.Create(c)) {
        assert.Equal(t, http.StatusCreated, rec.Code)
        var response domain.Transaction
        json.Unmarshal(rec.Body.Bytes(), &response)
        assert.Equal(t, validTransaction, response)
    }

    // Test invalid request
    req = httptest.NewRequest(http.MethodPost, "/transactions", bytes.NewReader([]byte("invalid")))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec = httptest.NewRecorder()
    c = e.NewContext(req, rec)

    if assert.NoError(t, handler.Create(c)) {
        assert.Equal(t, http.StatusBadRequest, rec.Code)
    }

    mockService.AssertExpectations(t)
}

func TestGetTransactionHandler(t *testing.T) {
    e := echo.New()
    mockService := new(MockService)
    handler := NewTransactionHandler(mockService)

    expectedTransaction := &domain.Transaction{
        ID:          1,
        FromAccount: "acc1",
        ToAccount:   "acc2",
        Amount:      100.0,
    }

    // Test successful get
    mockService.On("GetTransaction", 1).Return(expectedTransaction, nil)
    req := httptest.NewRequest(http.MethodGet, "/", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetPath("/transactions/:id")
    c.SetParamNames("id")
    c.SetParamValues("1")

    if assert.NoError(t, handler.GetByID(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
        var response domain.Transaction
        json.Unmarshal(rec.Body.Bytes(), &response)
        assert.Equal(t, *expectedTransaction, response)
    }

    // Test invalid ID
    req = httptest.NewRequest(http.MethodGet, "/", nil)
    rec = httptest.NewRecorder()
    c = e.NewContext(req, rec)
    c.SetPath("/transactions/:id")
    c.SetParamNames("id")
    c.SetParamValues("invalid")

    if assert.NoError(t, handler.GetByID(c)) {
        assert.Equal(t, http.StatusBadRequest, rec.Code)
    }

    // Test not found
    mockService.On("GetTransaction", 2).Return((*domain.Transaction)(nil), nil)
    req = httptest.NewRequest(http.MethodGet, "/", nil)
    rec = httptest.NewRecorder()
    c = e.NewContext(req, rec)
    c.SetPath("/transactions/:id")
    c.SetParamNames("id")
    c.SetParamValues("2")

    if assert.NoError(t, handler.GetByID(c)) {
        assert.Equal(t, http.StatusNotFound, rec.Code)
    }

    mockService.AssertExpectations(t)
}

func TestGetAllTransactionsHandler(t *testing.T) {
    e := echo.New()
    mockService := new(MockService)
    handler := NewTransactionHandler(mockService)

    expectedTransactions := []*domain.Transaction{
        {
            ID:          1,
            FromAccount: "acc1",
            ToAccount:   "acc2",
            Amount:      100.0,
        },
        {
            ID:          2,
            FromAccount: "acc3",
            ToAccount:   "acc4",
            Amount:      200.0,
        },
    }

    // Test successful get all
    mockService.On("GetAllTransactions").Return(expectedTransactions, nil)
    req := httptest.NewRequest(http.MethodGet, "/transactions", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)

    if assert.NoError(t, handler.GetAll(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
        var response []domain.Transaction
        json.Unmarshal(rec.Body.Bytes(), &response)
        assert.Equal(t, 2, len(response))
    }

    mockService.AssertExpectations(t)
}

func TestUpdateTransactionHandler(t *testing.T) {
    e := echo.New()
    mockService := new(MockService)
    handler := NewTransactionHandler(mockService)

    updateData := domain.Transaction{
        FromAccount: "acc1",
        ToAccount:   "acc2",
        Amount:      150.0,
    }

    // Test successful update
    mockService.On("UpdateTransaction", &domain.Transaction{
        ID:          1,
        FromAccount: "acc1",
        ToAccount:   "acc2",
        Amount:      150.0,
    }).Return(nil)

    jsonData, _ := json.Marshal(updateData)
    req := httptest.NewRequest(http.MethodPut, "/transactions/1", bytes.NewReader(jsonData))
    req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetPath("/transactions/:id")
    c.SetParamNames("id")
    c.SetParamValues("1")

    if assert.NoError(t, handler.Update(c)) {
        assert.Equal(t, http.StatusOK, rec.Code)
    }

    mockService.AssertExpectations(t)
}

func TestDeleteTransactionHandler(t *testing.T) {
    e := echo.New()
    mockService := new(MockService)
    handler := NewTransactionHandler(mockService)

    // Test successful delete
    mockService.On("DeleteTransaction", 1).Return(nil)
    req := httptest.NewRequest(http.MethodDelete, "/transactions/1", nil)
    rec := httptest.NewRecorder()
    c := e.NewContext(req, rec)
    c.SetPath("/transactions/:id")
    c.SetParamNames("id")
    c.SetParamValues("1")

    if assert.NoError(t, handler.Delete(c)) {
        assert.Equal(t, http.StatusNoContent, rec.Code)
    }

    mockService.AssertExpectations(t)
}