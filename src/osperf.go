package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
)

type osperf struct {
	object   string
	counter  string
	instance string
	value    int
	// nReads       int
	// nWrites      int
}

//Removes WhiteSpace
func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

//IO()
func OSperf(conn *sql.DB) (OSmetrics []osperf) {

	rows, _ := conn.Query("SELECT object_name, counter_name, instance_name, cntr_value  FROM sys.dm_os_performance_counters;") // Note: Ignoring errors for brevity
	// var name1 string
	var ob string
	var coun string
	var instan string
	var val int
	// var nr int
	// var nw int

	for rows.Next() {

		if err := rows.Scan(&ob, &coun, &instan, &val); err != nil {
			log.Fatal(err)
		}

		obLine := (standardizeSpaces(ob))
		counLine := (standardizeSpaces(coun))
		instanLine := (standardizeSpaces(instan))
		valLine := val

		os := osperf{
			object:   obLine,
			counter:  counLine,
			instance: instanLine,
			value:    valLine,
		}

		if strings.Contains(coun, "Buffer cache hit ratio base") {
			fmt.Println("Buffer cache hit ratio= ", val)
			OSmetrics = append(OSmetrics, os)
		}

		if strings.Contains(coun, "Page life expectancy") {
			fmt.Println("Page life expectancy= ", val)
			OSmetrics = append(OSmetrics, os)
		}

		if strings.Contains(coun, "Batch Requests") {
			fmt.Println("batch requests= ", val)
			OSmetrics = append(OSmetrics, os)
		}

		if strings.Contains(coun, "SQL Compilations") {
			fmt.Println("SQL Compilations= ", val)
			OSmetrics = append(OSmetrics, os)
		}

		if strings.Contains(coun, "Re-Compilations") {
			fmt.Println("SQL ReCompilations= ", val)
			OSmetrics = append(OSmetrics, os)
		}

		if strings.Contains(coun, "User Connections") {
			fmt.Println("User Connections= ", val)
			OSmetrics = append(OSmetrics, os)
		}

		if strings.Contains(coun, "Page Splits") {
			fmt.Println("Page Splits= ", val)
			OSmetrics = append(OSmetrics, os)
		}

		if strings.Contains(coun, "Processes blocked") {
			fmt.Println("Processes Blocked= ", val)
			OSmetrics = append(OSmetrics, os)
		}

		if strings.Contains(coun, "Checkpoint pages") {
			fmt.Println("Checkpoint Pages= ", val)
			OSmetrics = append(OSmetrics, os)
		}

		// fmt.Printf("***DatabaseName***: %s\n bytes: %d\n", dn, br)

		// os := osperf{
		// 	object:   obLine,
		// 	counter:  counLine,
		// 	instance: instanLine,
		// 	value:    valLine,
		// }

		// OSmetrics = append(OSmetrics, os)

	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return
}
