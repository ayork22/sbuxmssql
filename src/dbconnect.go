package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	debug         = flag.Bool("debug", false, "enable debugging")
	password      = flag.String("password", "test", "the database password")
	port     *int = flag.Int("port", 1433, "the database port")
	server        = flag.String("server", "DESKTOP-J0L7P9N\\SQLEXPRESS", "the database server")
	user          = flag.String("user", "test", "the database user")
)

var conn *sql.DB

//WinConnect
func DBconnect() (conn *sql.DB) {
	flag.Parse()

	if *debug {
		fmt.Printf(" password:%s\n", *password)
		fmt.Printf(" port:%d\n", *port)
		fmt.Printf(" server:%s\n", *server)
		fmt.Printf(" user:%s\n", *user)
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", *server, *user, *password, *port)
	if *debug {
		fmt.Printf(" connString:%s\n", connString)
	}
	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	// defer conn.Close()
	// defer stmt.Close()

	return
}
