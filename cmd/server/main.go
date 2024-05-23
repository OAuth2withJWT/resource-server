package main

import (
	"fmt"
	"log"

	"github.com/OAuth2withJWT/resource-server/db"
)

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Fatal("Failed to initialize database: ", err)
	}
	defer db.Close()

	fmt.Println("Hello, Resource Server")
}
