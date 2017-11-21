package main

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

type fileiotest struct {
	dbname       string
	bytesread    int
	byteswritten int
	sizeinbytes  int
	nReads       int
	nWrites      int
}

//IO()
func IO(conn *sql.DB) (IOmetrics []fileiotest) {
	rows, _ := conn.Query("SELECT d.name AS DatabaseName, SUM(a.num_of_bytes_read) AS BytesRead, SUM(a.num_of_bytes_written) AS BytesWritten, SUM(a.size_on_disk_bytes) AS SizeInBytes, SUM(a.num_of_reads) AS NumberOfReads, SUM(a.num_of_writes) AS NumberOfWrites FROM sys.databases d LEFT JOIN sys.dm_io_virtual_file_stats(NULL, NULL) a ON d.database_id = a.database_id GROUP BY d.name ORDER BY d.name") // Note: Ignoring errors for brevity
	// var name1 string

	var dn string
	var br int
	var bw int
	var sb int
	var nr int
	var nw int

	for rows.Next() {

		if err := rows.Scan(&dn, &br, &bw, &sb, &nr, &nw); err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("***DatabaseName***: %s\n bytes: %d\n", dn, br)

		io := fileiotest{
			dbname:       dn,
			bytesread:    br,
			byteswritten: bw,
			sizeinbytes:  sb,
			nReads:       nr,
			nWrites:      nw,
		}

		IOmetrics = append(IOmetrics, io)

	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return
}
