package main

import (
	"log"
	"net/http"

	// "database/sql"

	"library/handler"

	"github.com/gorilla/schema"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	var createTable = `
	CREATE TABLE IF NOT EXISTS category (
		id	serial,
		name text,
		status boolean,

		primary key (id)
	);`

	db, err := sqlx.Connect("postgres", "user=postgres password=password dbname=library sslmode=disable")
    if err != nil {
        log.Fatalln(err)
    }

	db.MustExec(createTable)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	r := handler.New(db, decoder)

	log.Println("Server starting...")
	if err:= http.ListenAndServe("127.0.0.1:3000", r); err != nil {
		log.Fatal(err)
	}
}