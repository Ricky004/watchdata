package watchdataexporter

import (
	"context"
	"fmt"
	"time"

	"github.com/Ricky004/watchdata/pkg/clickhousestore"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type watchdataExporter struct {
	dsn         string
	tlsInsecure bool
	logger      *zap.Logger
	ch     *clickhousestore.ClickHouseProvider
}

func newLogsExporter(cfg *Config, set exporter.Settings, ch *clickhousestore.ClickHouseProvider) (*watchdataExporter, error) {
	if cfg.DSN == "" {
		return nil, fmt.Errorf("DSN must be provided for watchdataExporter")
	}

	return &watchdataExporter{
		dsn:         cfg.DSN,
		tlsInsecure: cfg.TLSInsecure,
		logger:      set.Logger,
		ch:          ch,
	}, nil
}

// createLogsExporter is the factory function for the logs exporter.
func createLogsExporter(
	ctx context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Logs, error) {
	conf, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("unexpected config type: %T", cfg)
	}

	clickhouseCfg := clickhousestore.Config{
		Connection: clickhousestore.ConnectionConfig{
			DialTimeout: 5 * time.Second,
		},
	}

	chProvider, err := clickhousestore.NewClickHouseProvider(ctx, clickhouseCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to init ClickHouse: %w", err)
	}

	exp, err := newLogsExporter(conf, set, chProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to create watchdata logs exporter: %w", err)
	}

	return exp, nil
}

// Start is a lifecycle function for the exporter.
func (e *watchdataExporter) Start(ctx context.Context, host component.Host) error {
	e.logger.Info("Starting watchdataExporter with DSN", zap.String("dsn", e.dsn))
	return nil
}

// Shutdown is a lifecycle function for the exporter.
func (e *watchdataExporter) Shutdown(ctx context.Context) error {
	e.logger.Info("Stoping watchdataExporter with DSN", zap.String("dsn", e.dsn))
	return nil
}

// Capabilities returns the capabilities of the exporter.
func (e *watchdataExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

// ConsumeLogs is the method that receives log data.
func (e *watchdataExporter) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	records := convertToLogRecords(ld)
	err := e.ch.InsertLogs(ctx, records)
	if err != nil {
		return fmt.Errorf("failed to insert logs: %w", err)
	}
	return nil
}

// Compile-time check to ensure watchdataExporter implements exporter.Logs.
// If this line itself causes a compile error, it confirms the interface is not satisfied.
var _ exporter.Logs = (*watchdataExporter)(nil)
