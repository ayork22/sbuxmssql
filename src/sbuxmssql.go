package main

import (
	"fmt"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/metric"
	"github.com/newrelic/infra-integrations-sdk/sdk"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
	Duser string `default:"test" help:"File location to monitor"`
	Dpass string `default:"test" help:"File name being monitored"`
}

const (
	integrationName    = "com.newrelic.sbuxmssql"
	integrationVersion = "0.1.0"
)

var args argumentList

func populateInventory(inventory sdk.Inventory) error {
	// Insert here the logic of your integration to get the inventory data
	// Ex: inventory.SetItem("softwareVersion", "value", "1.0.1")
	// --
	return nil
}

func populateMetrics(ms *metric.MetricSet) error {

	//OPEN Databse Connection
	fmt.Println("***DUser***", args.Duser)
	fmt.Println("***DPass***", args.Dpass)
	fmt.Println("***BEFORE DBconnect***")
	var db = DBconnect(args.Duser, args.Dpass)
	fmt.Println("***Database Connected***")

	// IO Metrics
	var fileio = IO(db)
	for i := 0; i < len(fileio); i++ {

		// *****TESTING*****
		// ms.SetMetric("DatabaseName", fileio[i].dbname, metric.ATTRIBUTE)
		// *****TESTING*****
		ms.SetMetric("BytesRead_"+fileio[i].dbname, fileio[i].bytesread, metric.GAUGE)
		ms.SetMetric("BytesWritten_"+fileio[i].dbname, fileio[i].byteswritten, metric.GAUGE)
		ms.SetMetric("SizeInBytes_"+fileio[i].dbname, fileio[i].sizeinbytes, metric.GAUGE)
		ms.SetMetric("NumberReads_"+fileio[i].dbname, fileio[i].nReads, metric.GAUGE)
		ms.SetMetric("NumberWrites_"+fileio[i].dbname, fileio[i].nWrites, metric.GAUGE)
	}
	//
	// // Memory Metrics
	var mem = Memory(db)
	ms.SetMetric("BufferCache", mem.buffercache, metric.GAUGE)
	ms.SetMetric("PageLife", mem.pagelife, metric.GAUGE)

	// // Connections Metrics
	var connects = Connections(db)
	for i := 0; i < len(connects); i++ {
		ms.SetMetric("NumberConnections_"+connects[i].dbname, connects[i].nConnections, metric.GAUGE)
		ms.SetMetric("NumberReadsConnections_"+connects[i].dbname, connects[i].nReads, metric.GAUGE)
		ms.SetMetric("NumberWritesConnections_"+connects[i].dbname, connects[i].nWrites, metric.GAUGE)
	}

	// CLOSE Databse Connection
	defer db.Close()
	return nil
}

func main() {
	integration, err := sdk.NewIntegration(integrationName, integrationVersion, &args)
	fatalIfErr(err)

	fmt.Println("***PLUGIN STARTED***")

	if args.All || args.Inventory {
		fatalIfErr(populateInventory(integration.Inventory))
	}

	if args.All || args.Metrics {
		ms := integration.NewMetricSet("MSSQL")
		fatalIfErr(populateMetrics(ms))
	}
	fatalIfErr(integration.Publish())
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
