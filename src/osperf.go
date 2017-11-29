package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/newrelic/infra-integrations-sdk/metric"
	"github.com/newrelic/infra-integrations-sdk/sdk"
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
func OSperf(conn *sql.DB, test *sdk.Integration) {

	f, err := os.OpenFile("C:\\Users\\ayork\\Desktop\\test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	rows, _ := conn.Query("SELECT object_name, counter_name, instance_name, cntr_value  FROM sys.dm_os_performance_counters;") // Note: Ignoring errors for brevity
	defer rows.Close()
	var ob string
	var coun string
	var instan string
	var val int
	// var nr int
	// var nw int
	ms := test.NewMetricSet("OSperf")
	for rows.Next() {

		if err := rows.Scan(&ob, &coun, &instan, &val); err != nil {
			log.Fatal(err)
		}

		// obLine := (standardizeSpaces(ob))
		counLine := (standardizeSpaces(coun))
		// instanLine := (standardizeSpaces(instan))
		valLine := val

		if counLine == "Buffer cache hit ratio" {
			fmt.Println("Buffer cache hit ratio= ", val)
			setMertric(ms, counLine, valLine, metric.GAUGE)
		} else if counLine == "Buffer cache hit ratio base" {
			fmt.Println("Buffer cache hit ratio base= ", val)
			setMertric(ms, counLine, valLine, metric.GAUGE)
		} else if strings.Contains(coun, "Page life expectancy") {
			fmt.Println("Page life expectancy= ", val)
			setMertric(ms, counLine, valLine, metric.GAUGE)
		} else if strings.Contains(coun, "Batch Requests") {
			fmt.Println("batch requests= ", val)
			setMertric(ms, counLine, valLine, metric.GAUGE)
		} else if strings.Contains(coun, "SQL Compilations") {
			fmt.Println("SQL Compilations= ", val)
			setMertric(ms, counLine, valLine, metric.GAUGE)
		} else if strings.Contains(coun, "Re-Compilations") {
			fmt.Println("SQL ReCompilations= ", val)
			setMertric(ms, counLine, valLine, metric.GAUGE)
		} else if strings.Contains(coun, "User Connections") {
			fmt.Println("User Connections= ", val)
			setMertric(ms, counLine, valLine, metric.GAUGE)
		} else if strings.Contains(coun, "Page Splits") {
			fmt.Println("Page Splits= ", val)
			setMertric(ms, counLine, valLine, metric.GAUGE)
		} else if strings.Contains(coun, "Processes blocked") {
			fmt.Println("Processes Blocked= ", val)
			setMertric(ms, counLine, valLine, metric.GAUGE)
		} else if strings.Contains(coun, "Checkpoint pages") {
			fmt.Println("Checkpoint Pages= ", val)
			setMertric(ms, counLine, valLine, metric.GAUGE)
		}

	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

}
