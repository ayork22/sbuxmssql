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

	stmt, err := conn.Prepare("select 1, 'abc'")
	if err != nil {
		log.Fatal("Prepare failed:", err.Error())
	}
	defer stmt.Close()

	row := stmt.QueryRow()
	var somenumber int64
	var somechars string
	err = row.Scan(&somenumber, &somechars)
	if err != nil {
		log.Fatal("Scan failed:", err.Error())
	}

	rows1, _ := conn.Query("SELECT (SELECT cntr_value * 100.00 FROM sys.dm_os_performance_counters WHERE counter_name = 'Buffer cache hit ratio') / (SELECT cntr_value FROM sys.dm_os_performance_counters WHERE counter_name = 'Buffer cache hit ratio base') AS BufferCacheHitRatio, (SELECT cntr_value FROM sys.dm_os_performance_counters WHERE counter_name = 'Page life expectancy' AND RTRIM([object_name]) LIKE '%:Buffer Manager') AS PageLife") // Note: Ignoring errors for brevity

	var bc float64
	var pl int

	for rows1.Next() {
		if err := rows1.Scan(&bc, &pl); err != nil {
			log.Fatal(err)
		}
		fmt.Println(bc, pl)

	}
	fmt.Printf("somenumber:%d\n", somenumber)
	fmt.Printf("somechars:%s\n", somechars)

	fmt.Printf("bye\n")

	return
}
