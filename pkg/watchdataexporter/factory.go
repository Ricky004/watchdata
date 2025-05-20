package watchdataexporter

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
)

var (
	typStr = component.MustNewType("watchdataexporter")
)

func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		typStr,
		CreateDefaultConfig,
		exporter.WithLogs(createLogsExporter, component.StabilityLevelAlpha),
	)
}

func CreateDefaultConfig() component.Config {
	return &Config{
		Endpoint: "http://host.docker.internal:4320",
		APIKey:   "",
	}
}
