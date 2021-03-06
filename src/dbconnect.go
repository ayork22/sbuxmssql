package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
)

// var (
// 	debug         = flag.Bool("debug", true, "enable debugging")
// 	password      = flag.String("password", "test", "the database password")
// 	port     *int = flag.Int("port", 1433, "the database port")
// 	server        = flag.String("server", "DESKTOP-J0L7P9N\\SQLEXPRESS", "the database server")
// 	dbuser        = flag.String("dbuser", "test", "the database user")
// )

var (
	debug         = flag.Bool("debug", true, "enable debugging")
	password      = flag.String("password", "c0y0te#22", "the database password")
	port     *int = flag.Int("port", 1433, "the database port")
	server        = flag.String("server", "localhost", "the database server")
	dbuser        = flag.String("dbuser", "sa", "the database user")
)

// var conn *sql.DB

//WinConnect
func DBconnect(du, dp, dserver string) (conn *sql.DB) {
	flag.Parse()

	if *debug {
		fmt.Printf(" password:%s\n", *password)
		fmt.Printf(" port:%d\n", *port)
		fmt.Printf(" server:%s\n", *server)
		fmt.Printf(" dbuser:%s\n", *dbuser)
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	// Windows the DB Server name is needed as well
	// hostname = hostname + dserver
	fmt.Println("OS Hostname:", hostname)

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", hostname, du, dp, *port)
	if *debug {
		fmt.Printf("connString:%s\n", connString)
	}
	conn, err = sql.Open("mssql", connString)
	// if err != nil {
	// 	log.Fatal("Open connection failed:", err.Error())
	// }

	if err != nil {
		fmt.Println("TEST", err)
		os.Exit(1)
	}

	err = conn.Ping()

	if err != nil {
		fmt.Println("Can't Connect", err)
		os.Exit(1)
	}

	// if no error. Ping is successful
	fmt.Println("Ping to database successful, connection is still alive")
	// conn.Close()

	// this part should fail

	err = conn.Ping()

	if err != nil {
		fmt.Println("ErrorMessage: ", err)
		os.Exit(1)
	}
	// if no error. Ping is successful
	fmt.Println("Ping to database successful, connection is still alive")

	return
}
