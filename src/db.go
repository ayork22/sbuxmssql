package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	debug    = flag.Bool("debug", true, "enable debugging")
	dbuser   = flag.String("dbuser", "sa", "the database user")
	password = flag.String("password", "c0y0te#22", "the database password")
	port     = flag.Int("port", 1433, "the database port")
	server   = flag.String("server", "localhost", "the database server")
)

// SELECT (SELECT cntr_value * 100.00 FROM sys.dm_os_performance_counters WHERE counter_name = 'Buffer cache hit ratio') / (SELECT cntr_value FROM sys.dm_os_performance_counters WHERE counter_name = 'Buffer cache hit ratio base') AS BufferCacheHitRatio, (SELECT cntr_value FROM sys.dm_os_performance_counters WHERE counter_name = 'Page life expectancy' AND RTRIM([object_name]) LIKE '%:Buffer Manager') AS PageLife
type memory struct {
	buffercache float64
	pagelife    int
}

// SELECT d.Name AS DatabaseName, COUNT(c.connection_id) AS NumberOfConnections, ISNULL(SUM(c.num_reads), 0) AS NumberOfReads, ISNULL(SUM(c.num_writes), 0) AS NumberOfWrites FROM sys.databases d LEFT JOIN sys.sysprocesses s ON s.dbid = d.database_id LEFT JOIN sys.dm_exec_connections c ON c.session_id = s.spid WHERE (s.spid IS NULL OR c.session_id >= 51) GROUP BY d.Name
type connections struct {
	dbname       string
	nConnections int
	nReads       int
	nWrites      int
}

// SELECT d.name AS DatabaseName, SUM(a.num_of_bytes_read) AS BytesRead, SUM(a.num_of_bytes_written) AS BytesWritten, SUM(a.size_on_disk_bytes) AS SizeInBytes, SUM(a.num_of_reads) AS NumberOfReads, SUM(a.num_of_writes) AS NumberOfWrites FROM sys.databases d LEFT JOIN sys.dm_io_virtual_file_stats(NULL, NULL) a ON d.database_id = a.database_id GROUP BY d.name ORDER BY d.name
type fileio struct {
	dbname       string
	bytesread    int
	byteswritten int
	sizeinbytes  int
	nReads       int
	nWrites      int
}

//Calldb test
// func Calldb() (IOmetrics []fileio, MEMmetrics []memory) {
func Calldb() (IOmetrics []fileio) {
	flag.Parse()

	if *debug {
		fmt.Printf(" dbuser:%s\n", *dbuser)
		fmt.Printf(" password:%s\n", *password)
		fmt.Printf(" port:%d\n", *port)
		fmt.Printf(" server:%s\n", *server)
	}

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d", *server, *dbuser, *password, *port)
	if *debug {
		fmt.Printf(" connString:%s\n", connString)
	}
	conn, err := sql.Open("mssql", connString)
	fmt.Printf("***Connected***\n")
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}

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
		// fmt.Printf("%v\n", sqlversion)
	}

	// rows, _ := conn.Query("SELECT d.name AS DatabaseName, SUM(a.num_of_bytes_read) AS BytesRead FROM sys.databases d LEFT JOIN sys.dm_io_virtual_file_stats(NULL, NULL) a ON d.database_id = a.database_id GROUP BY d.name ORDER BY d.name") // Note: Ignoring errors for brevity
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

		io := fileio{
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

	// rows, err := conn.Query("SELECT d.name AS DatabaseName, SUM(a.num_of_bytes_read) AS BytesRead FROM sys.databases d LEFT JOIN sys.dm_io_virtual_file_stats(NULL, NULL) a ON d.database_id = a.database_id GROUP BY d.name ORDER BY d.name")
	// rows1, err := conn.Query("SELECT d.name AS DatabaseName, SUM(a.num_of_bytes_read) AS BytesRead, SUM(a.num_of_bytes_written) AS BytesWritten, SUM(a.size_on_disk_bytes)	AS SizeInBytes, SUM(a.num_of_reads)	AS NumberOfReads, SUM(a.num_of_writes) AS NumberOfWrites FROM sys.databases d LEFT JOIN sys.dm_io_virtual_file_stats(NULL, NULL) a ON d.database_id = a.database_id GROUP BY d.name ORDER BY d.name")

	rows2, _ := conn.Query("SELECT (SELECT cntr_value * 100.00 FROM sys.dm_os_performance_counters WHERE counter_name = 'Buffer cache hit ratio') / (SELECT cntr_value FROM sys.dm_os_performance_counters WHERE counter_name = 'Buffer cache hit ratio base') AS BufferCacheHitRatio, (SELECT cntr_value FROM sys.dm_os_performance_counters WHERE counter_name = 'Page life expectancy' AND RTRIM([object_name]) LIKE '%:Buffer Manager') AS PageLife") // Note: Ignoring errors for brevity
	// var name1 string

	var bc float64
	var pl int

	for rows2.Next() {

		if err := rows2.Scan(&bc, &pl); err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("***DatabaseName***: %s\n bytes: %d\n", dn, br)

		mem := memory{
			buffercache: bc,
			pagelife:    pl,
		}
		var MEMmetrics []memory
		MEMmetrics = append(MEMmetrics, mem)

		// fmt.Println(mem.pagelife, "\n")
		// fmt.Println(mem.buffercache)

	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	fmt.Printf("***bye***\n")
	return
}
