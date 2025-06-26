package builderconfig

import (
	"github.com/Ricky004/watchdata/pkg/types/otelpipelinetypes"
)

var ReceiverModules = map[string]otelpipelinetypes.ModuleEntry{
	"otlp": {
		Gomod: "go.opentelemetry.io/collector/receiver/otlpreceiver v0.128.0",
	},
	"filelog": {
		Gomod: "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/filelogreceiver v0.127.0",
	},
}

var ProcessorModules = map[string]otelpipelinetypes.ModuleEntry{
	"batch": {
		Gomod: "go.opentelemetry.io/collector/processor/batchprocessor v0.127.0",
	},
}

var ExporterModules = map[string]otelpipelinetypes.ModuleEntry{
	"watchdataexporter": {
		Gomod:  "github.com/Ricky004/watchdata v0.0.7-0.20250610145142-6cee7eaf3324",
		Import: "github.com/Ricky004/watchdata/pkg/watchdataexporter",
	},
}

func SyncBuilderConfig(path string, receivers, processors, exporters []string) error {
	cfg, err := LoadBuilderConfig(path)
	if err != nil {
		return err
	}

	for _, name := range receivers {
		if mod, ok := ReceiverModules[name]; ok {
			AddModule(&cfg.Receivers, mod.Gomod, mod.Import)
		}
	}

	for _, name := range processors {
		if mod, ok := ProcessorModules[name]; ok {
			AddModule(&cfg.Processors, mod.Gomod, mod.Import)
		}
	}

	for _, name := range exporters {
		if mod, ok := ExporterModules[name]; ok {
			AddModule(&cfg.Exporters, mod.Gomod, mod.Import)
		}
	}

	return SaveBuilderConfig(path, cfg)
}
