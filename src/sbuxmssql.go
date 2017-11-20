package main

import (
	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/metric"
	"github.com/newrelic/infra-integrations-sdk/sdk"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
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
	// Ex: ms.SetMetric("requestsPerSecond", 10, metric.GAUGE)

	var fileio = Calldb()
	for i := 0; i < len(fileio); i++ {
		ms.SetMetric("BytesRead_"+fileio[i].dbname, fileio[i].bytesread, metric.GAUGE)
		ms.SetMetric("BytesWritten_"+fileio[i].dbname, fileio[i].byteswritten, metric.GAUGE)
		ms.SetMetric("SizeInBytes_"+fileio[i].dbname, fileio[i].sizeinbytes, metric.GAUGE)
		ms.SetMetric("NumberReads_"+fileio[i].dbname, fileio[i].nReads, metric.GAUGE)
		ms.SetMetric("NumberWrites_"+fileio[i].dbname, fileio[i].nWrites, metric.GAUGE)
	}
	return nil
}

func main() {
	integration, err := sdk.NewIntegration(integrationName, integrationVersion, &args)
	fatalIfErr(err)

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
