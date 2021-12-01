package main

import (
	"log"
	"net/http"

	// "database/sql"

	"library/handler"

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	var createTable = `
	CREATE TABLE IF NOT EXISTS categories (
		id	serial,
		name text,
		status boolean,

		primary key (id)
	);
	
	CREATE TABLE IF NOT EXISTS books (
		id	serial,
		category_id integer,
		book_name text,
		author_name text,
		details text,
		image text,
		status boolean,

		primary Key (id)
	);
	
	CREATE TABLE IF NOT EXISTS bookings (
		id	serial,
		user_id integer,
		book_id integer,
		start_time timestamp,
		end_time timestamp,

		primary Key (id)
	);
	
	CREATE TABLE IF NOT EXISTS users (
		id	serial,
		first_name text,
		last_name text,
		email text,
		password text,
		is_verified boolean,

		primary Key (id)
	);`

	db, err := sqlx.Connect("postgres", "user=postgres password=password dbname=library sslmode=disable")
    if err != nil {
        log.Fatalln(err)
    }

	db.MustExec(createTable)
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	store := sessions.NewCookieStore([]byte("jsowjpw38eowj4ur82jmaole0uehqpl"))
	r := handler.New(db, decoder, store)

	log.Println("Server starting...")
	if err:= http.ListenAndServe("127.0.0.1:3000", r); err != nil {
		log.Fatal(err)
	}
}