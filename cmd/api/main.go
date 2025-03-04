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

    e.POST("/transactions", handler.Create)
    e.GET("/transactions/:id", handler.GetByID)
    e.GET("/transactions", handler.GetAll)
    e.PUT("/transactions/:id", handler.Update)
    e.DELETE("/transactions/:id", handler.Delete)

    e.Logger.Fatal(e.Start(":8080"))
} 