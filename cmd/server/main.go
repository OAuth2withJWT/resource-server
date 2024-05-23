package main

import (
	"fmt"
	"log"

	"github.com/OAuth2withJWT/resource-server/app"
	"github.com/OAuth2withJWT/resource-server/app/postgres"
	"github.com/OAuth2withJWT/resource-server/db"
	"github.com/OAuth2withJWT/resource-server/server"
)

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	defer db.Close()

	cardRepository := postgres.NewCardRepository(db)
	transactionRepository := postgres.NewTransactionRepository(db)

	app := app.Application{
		CardService:        app.NewCardService(cardRepository),
		TransactionService: app.NewTransactionService(transactionRepository),
	}
	s := server.New(&app)
	log.Fatal(s.Run())

	fmt.Println("Hello, Resource Server")
}
