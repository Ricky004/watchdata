package otelexporter

import (
	"go.opentelemetry.io/collector/component"
)

// Config defines configuration for the WatchData exporter.
type Config struct {
	ID       component.ID `mapstructure:"-"`
	Endpoint string       `mapstructure:"endpoint"`
	APIKey   string       `mapstructure:"api_key"`
}

// GetID return exporter ID.
func (cfg *Config) GetID() component.ID {
	return cfg.ID
}