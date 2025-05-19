package otelexporter

import (

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
)

const typeStr string = "watchdata"

func NewFactory() exporter.Factory {
	typ := component.MustNewType(typeStr)
	return exporter.NewFactory(
		typ,
		CreateDefaultConfig,
		exporter.WithLogs(createLogsExporter, component.StabilityLevelBeta),
	)
}

func CreateDefaultConfig() component.Config {
	typ := component.MustNewType(typeStr)
	return &Config{
		ID:       component.NewID(typ),
		Endpoint: "http://host.docker.internal:4320",
		APIKey:   "",
	}
}
