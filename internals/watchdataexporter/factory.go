package watchdataexporter

import (
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
)

const typeStr string = "watchdata"

func NewFactory() component.Factory {
	typ, err := component.NewType(typeStr)
	if err != nil {
		fmt.Println("Type is invalid", err)
	}
	return exporter.NewFactory(
		typ,
		CreateDefaultConfig,
		exporter.WithLogs(createLogsExporter, component.StabilityLevelBeta),
	)
}

func CreateDefaultConfig() component.Config {
	typ, _ := component.NewType(typeStr)
	return &Config{
		ID:       component.NewID(typ),
		Endpoint: "http://host.docker.internal:4320",
		APIKey:   "",
	}
}
