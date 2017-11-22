package main

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

// SELECT d.Name AS DatabaseName, COUNT(c.connection_id) AS NumberOfConnections, ISNULL(SUM(c.num_reads), 0) AS NumberOfReads, ISNULL(SUM(c.num_writes), 0) AS NumberOfWrites FROM sys.databases d LEFT JOIN sys.sysprocesses s ON s.dbid = d.database_id LEFT JOIN sys.dm_exec_connections c ON c.session_id = s.spid WHERE (s.spid IS NULL OR c.session_id >= 51) GROUP BY d.Name
type connections struct {
	dbname       string
	nConnections int
	nReads       int
	nWrites      int
}

//IO()
func Connections(conn *sql.DB) (CONNECTIONmetrics []connections) {
	rows, _ := conn.Query("SELECT d.Name AS DatabaseName, COUNT(c.connection_id) AS NumberOfConnections, ISNULL(SUM(c.num_reads), 0) AS NumberOfReads, ISNULL(SUM(c.num_writes), 0) AS NumberOfWrites FROM sys.databases d LEFT JOIN sys.sysprocesses s ON s.dbid = d.database_id LEFT JOIN sys.dm_exec_connections c ON c.session_id = s.spid WHERE (s.spid IS NULL OR c.session_id >= 51) GROUP BY d.Name") // Note: Ignoring errors for brevity
	// var name1 string

	var dn string
	var nc int
	var nr int
	var nw int

	for rows.Next() {

		if err := rows.Scan(&dn, &nc, &nr, &nw); err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("***DatabaseName***: %s\n bytes: %d\n", dn, br)

		c := connections{
			dbname:       dn,
			nConnections: nc,
			nReads:       nr,
			nWrites:      nw,
		}

		CONNECTIONmetrics = append(CONNECTIONmetrics, c)

	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return
}
