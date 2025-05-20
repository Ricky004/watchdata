package watchdataexporter

// Config defines configuration for the WatchData exporter.
type Config struct {
	Endpoint string `mapstructure:"endpoint"`
	APIKey   string `mapstructure:"api_key"`
}
