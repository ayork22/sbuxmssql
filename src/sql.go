package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

func SQL(conn *sql.DB) {
	var (
		sqlversion string
	)
	rows1, err := conn.Query("SELECT name FROM master..sysdatabases;")
	if err != nil {
		log.Fatal(err)
	}
	for rows1.Next() {
		err := rows1.Scan(&sqlversion)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%v\n SQL.GO TEST HERE WORKING!!!!!", sqlversion)
	}
}
