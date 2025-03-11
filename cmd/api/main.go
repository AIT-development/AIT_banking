package main

import (
    "bank-transactions/internal/handler"
    "bank-transactions/internal/repository"
    "bank-transactions/internal/service"
    "database/sql"
    "log"
    
    "github.com/labstack/echo/v4"
    _ "github.com/lib/pq"
)

func main() {
    db, err := sql.Open("postgres", "postgres://user:password@localhost/bank?sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    repo := repository.NewTransactionRepository(db)
    service := service.NewTransactionService(repo)
    handler := handler.NewTransactionHandler(service)

    e := echo.New()

    // @title Create Transaction
    // @description Creates a new transaction
    e.POST("/transactions", handler.Create)

    // @title Get Transaction by ID
    // @description Retrieves a specific transaction by its ID
    e.GET("/transactions/:id", handler.GetByID)

    // @title Get All Transactions
    // @description Retrieves all transactions
    e.GET("/transactions", handler.GetAll)

    // @title Update Transaction
    // @description Updates an existing transaction
    e.PUT("/transactions/:id", handler.Update)

    // @title Delete Transaction
    // @description Deletes a transaction
    e.DELETE("/transactions/:id", handler.Delete)

    e.Logger.Fatal(e.Start(":8080"))
} 