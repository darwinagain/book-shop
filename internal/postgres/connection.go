package postgres

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func OpenConnection() *sql.DB {
	db, err := sql.Open("postgres", "postgres://pguser:pgpass@localhost:5432/book-shop?sslmode=disable")
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	return db
}
