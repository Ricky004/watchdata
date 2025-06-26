package exporter

var ExporterTemplates = map[string]interface{}{
	"watchdataexporter": map[string]interface{}{
		"dsn": "tcp://clickhouse:9000/default?username=default&password=pass",
		"insecure": true,
	},
}