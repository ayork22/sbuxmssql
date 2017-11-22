package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var db *sql.DB

func DBconnect() (db *sql.DB) {
	var err error

	fmt.Println("Database Connect Function Called")

	db, err = sql.Open("mssql", "server=localhost;user id=sa;password=c0y0te#22;port=1433")

	if err != nil {
		fmt.Println("Database Connection Failed")
		log.Fatal(err)

	}

	// Other setup-related activities
	return
}
