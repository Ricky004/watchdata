package otelpipelinetypes

type OTelConfig struct {
	Receivers  map[string]interface{} `yaml:"receivers"`
	Processors map[string]interface{} `yaml:"processors"`
	Exporters  map[string]interface{} `yaml:"exporters"`
	Service    ServiceConfig           `yaml:"service"`
}

type ServiceConfig struct {
	Pipelines map[string]Pipeline `yaml:"pipelines"`
}

type Pipeline struct {
	Receivers  []string `yaml:"receivers"`
	Processors []string `yaml:"processors"`
	Exporters  []string `yaml:"exporters"`
}
