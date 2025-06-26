package otelpipelinetypes

type BuilderConfigs struct {
	Dist struct {
		Name        string `yaml:"name"`
		Description string `yaml:"description"`
		OutputPath  string `yaml:"output_path"`
	} `yaml:"dist"`

	Receivers  []ModuleEntry `yaml:"receivers"`
	Processors []ModuleEntry `yaml:"processors"`
	Exporters  []ModuleEntry `yaml:"exporters"`
}

type ModuleEntry struct {
	Gomod  string `yaml:"gomod"`
	Import string `yaml:"import,omitempty"`
}
