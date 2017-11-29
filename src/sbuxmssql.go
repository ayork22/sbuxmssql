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
	Duser   string `default:"sa" help:"Database UserName"`
	Dpass   string `default:"c0y0te#22" help:"Database Password"`
	Dserver string `default:"SQLEXPRESS" help:"Database Server Name"`
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

// func populateMetrics(ms *metric.MetricSet) error {
// Instead of just creating 1 new Metric Set I'm creating a new MetricSet for each Row I loop through
func populateMetrics(integration *sdk.Integration) error {

	//OPEN Databse Connection
	fmt.Println("***DUser***", args.Duser)
	fmt.Println("***DPass***", args.Dpass)
	fmt.Println("***BEFORE DBconnect***")
	var db = DBconnect(args.Duser, args.Dpass, "\\"+args.Dserver)
	fmt.Println("***Database Connected***")

	//OS PERF
	OSperf(db, integration)
	// var osPerf = OSperf(db, integration)
	// for i := 0; i < len(osPerf); i++ {
	// 	ms := integration.NewMetricSet("OSperfTest")
	// 	setMertric(ms, "ObjectName", osPerf[i].object, metric.GAUGE)
	// 	setMertric(ms, "CounterName", osPerf[i].counter, metric.GAUGE)
	// 	setMertric(ms, "DBName", osPerf[i].instance, metric.GAUGE)
	// 	setMertric(ms, "Value", osPerf[i].value, metric.GAUGE)
	//
	// 	// setMertric(ms, "NumberReads", fileio[i].nReads, metric.GAUGE)
	// 	// setMertric(ms, "NumberWrites", fileio[i].nWrites, metric.GAUGE)
	// }

	// IO Metrics
	var fileio = IO(db)
	for i := 0; i < len(fileio); i++ {
		ms := integration.NewMetricSet("MSSQL")
		setMertric(ms, "DatabaseName", fileio[i].dbname, metric.ATTRIBUTE)
		setMertric(ms, "BytesRead", fileio[i].bytesread, metric.GAUGE)
		setMertric(ms, "BytesWritten", fileio[i].byteswritten, metric.GAUGE)
		setMertric(ms, "SizeInBytes", fileio[i].sizeinbytes, metric.GAUGE)
		setMertric(ms, "NumberReads", fileio[i].nReads, metric.GAUGE)
		setMertric(ms, "NumberWrites", fileio[i].nWrites, metric.GAUGE)
		// ****TESTING*****
		// ms.SetMetric("DatabaseName", fileio[i].dbname, metric.GAUGE)
		// ms.SetMetric("BytesRead_", fileio[i].bytesread, metric.GAUGE)
		// ms.SetMetric("BytesWritten_", fileio[i].byteswritten, metric.GAUGE)
		// ms.SetMetric("SizeInBytes_", fileio[i].sizeinbytes, metric.GAUGE)
		// ms.SetMetric("NumberReads_", fileio[i].nReads, metric.GAUGE)
		// ms.SetMetric("NumberWrites_", fileio[i].nWrites, metric.GAUGE)
	}
	//
	// // Memory Metrics
	var mem = Memory(db)
	ms := integration.NewMetricSet("MSSQL")
	setMertric(ms, "BufferCache", mem.buffercache, metric.GAUGE)
	setMertric(ms, "PageLife", mem.pagelife, metric.GAUGE)
	// ms.SetMetric("BufferCache", mem.buffercache, metric.GAUGE)
	// ms.SetMetric("PageLife", mem.pagelife, metric.GAUGE)

	// // Connections Metrics
	var connects = Connections(db)
	for i := 0; i < len(connects); i++ {
		ms := integration.NewMetricSet("MSSQL")
		setMertric(ms, "DatabaseName", connects[i].dbname, metric.ATTRIBUTE)
		setMertric(ms, "NumberConnections", connects[i].nConnections, metric.GAUGE)
		setMertric(ms, "NumberReadsConnections", connects[i].nReads, metric.GAUGE)
		setMertric(ms, "NumberWritesConnections", connects[i].nWrites, metric.GAUGE)

		// ms.SetMetric("DatabaseName", connects[i].dbname, metric.GAUGE)
		// ms.SetMetric("NumberConnections", connects[i].nConnections, metric.GAUGE)
		// ms.SetMetric("NumberReadsConnections", connects[i].nReads, metric.GAUGE)
		// ms.SetMetric("NumberWritesConnections", connects[i].nWrites, metric.GAUGE)
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
		// ms := integration.NewMetricSet("MSSQL")
		fatalIfErr(populateMetrics(integration))
	}
	fatalIfErr(integration.Publish())
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func setMertric(metricSet *metric.MetricSet, name string, value interface{}, sourceType metric.SourceType) error {
	err := metricSet.SetMetric(name, value, sourceType)
	if err != nil {
		log.Warn("Failed setting value. name=%s  value=%s Error= %s", name, value, err)
	}
	return err
}
