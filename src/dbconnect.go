package main

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var db *sql.DB

func DBconnect() (db *sql.DB) {
	var err error

	db, err = sql.Open("mssql", "server=localhost;user id=sa;password=c0y0te#22;port=1433")
	if err != nil {
		log.Fatal(err)
	}

	// Other setup-related activities
	return
}
