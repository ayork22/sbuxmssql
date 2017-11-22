package main

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

type memory2 struct {
	buffercache float64
	pagelife    int
}

//Memory()
func Memory(conn *sql.DB) (MEMtest memory2) {
	rows, _ := conn.Query("SELECT (SELECT cntr_value * 100.00 FROM sys.dm_os_performance_counters WHERE counter_name = 'Buffer cache hit ratio') / (SELECT cntr_value FROM sys.dm_os_performance_counters WHERE counter_name = 'Buffer cache hit ratio base') AS BufferCacheHitRatio, (SELECT cntr_value FROM sys.dm_os_performance_counters WHERE counter_name = 'Page life expectancy' AND RTRIM([object_name]) LIKE '%:Buffer Manager') AS PageLife") // Note: Ignoring errors for brevity
	defer rows.Close()
	var bc float64
	var pl int

	for rows.Next() {
		if err := rows.Scan(&bc, &pl); err != nil {
			log.Fatal(err)
		}

		MEMtest = memory2{
			buffercache: bc,
			pagelife:    pl,
		}
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return
}
