package watchdataexporter

import (
	"go.opentelemetry.io/collector/component"
)

// Config defines configuration for the WatchData exporter.
type Config struct {
	component.Config `mapstructure:",squash"`
	Endpoint         string `mapstructure:"endpoint"`
	TLSInsecure      bool   `mapstructure:"insecure"`
}
