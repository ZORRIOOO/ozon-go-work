package main

import (
	"database/sql"
	"github.com/pressly/goose/v3"
	"log"
)

const connectionUrl = "postgres://user:password@localhost:5432/homework"

func main() {
	db, err := sql.Open("postgres", connectionUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("Ошибка применения миграций: %v", err)
	}
}
