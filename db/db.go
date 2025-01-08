package db

import (
    "database/sql"
    "log"

    _ "github.com/lib/pq"
)

func Connect() *sql.DB {
    dsn := "user=postgres password=123 dbname=db1 sslmode=disable"
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal(err)
    } else {
       log.Println("DB Connect Success")
    }
    return db
}
