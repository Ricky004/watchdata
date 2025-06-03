package watchdataexporter

import (
	"go.opentelemetry.io/collector/component"
)

// Config defines configuration for the WatchData exporter.
type Config struct {
	component.Config `mapstructure:",squash"`
	DSN              string `mapstructure:"dsn"`
	TLSInsecure      bool   `mapstructure:"insecure"`
}
