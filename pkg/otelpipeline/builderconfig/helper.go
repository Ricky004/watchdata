package builderconfig

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/Ricky004/watchdata/pkg/types/otelpipelinetypes"
)

func LoadBuilderConfig(path string) (*otelpipelinetypes.BuilderConfigs, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg otelpipelinetypes.BuilderConfigs
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func SaveBuilderConfig(path string, cfg *otelpipelinetypes.BuilderConfigs) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Better use AddModule
func AddModule(entries *[]otelpipelinetypes.ModuleEntry, gomod string, importpath ...string) {
	for _, e := range *entries {
		if e.Gomod == gomod {
			return
		}
	}

	entry := otelpipelinetypes.ModuleEntry{Gomod: gomod}
	if len(importpath) > 0 {
		entry.Import = importpath[0]
	}
	*entries = append(*entries, entry)
}

func AddReceiver(cfg *otelpipelinetypes.BuilderConfigs, gomod string) {
	for _, r := range cfg.Receivers {
		if r.Gomod == gomod {
			return
		}
	}

	cfg.Receivers = append(cfg.Receivers, otelpipelinetypes.ModuleEntry{Gomod: gomod})
}

func AddProcessor(cfg *otelpipelinetypes.BuilderConfigs, gomod string) {
	for _, r := range cfg.Processors {
		if r.Gomod == gomod {
			return
		}
	}

	cfg.Processors = append(cfg.Processors, otelpipelinetypes.ModuleEntry{Gomod: gomod})
}

func AddExporter(cfg *otelpipelinetypes.BuilderConfigs, gomod, importPath string) {
	for _, r := range cfg.Exporters {
		if r.Gomod == gomod {
			return
		}
	}
	cfg.Exporters = append(cfg.Exporters, otelpipelinetypes.ModuleEntry{
		Gomod:  gomod,
		Import: importPath,
	})
}

