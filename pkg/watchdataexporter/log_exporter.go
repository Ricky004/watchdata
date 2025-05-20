package watchdataexporter

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/plog/plogotlp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type watchdataExporter struct {
	endpoint    string
	tlsInsecure bool
	logger      *zap.Logger

	conn   *grpc.ClientConn
	client plogotlp.GRPCClient
	cancel context.CancelFunc
}

func newLogsExporter(cfg *Config, set exporter.Settings) (*watchdataExporter, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("endpoint must be provided for watchdataExporter")
	}

	return &watchdataExporter{
		endpoint:    cfg.Endpoint,
		tlsInsecure: cfg.TLSInsecure,
		logger:      set.Logger,
	}, nil
}

// createLogsExporter is the factory function for the logs exporter.
func createLogsExporter(
	_ context.Context,
	set exporter.Settings,
	cfg component.Config,
) (exporter.Logs, error) {
	conf, ok := cfg.(*Config)
	if !ok {
		return nil, fmt.Errorf("unexpected config type: %T", cfg)
	}

	exp, err := newLogsExporter(conf, set)
	if err != nil {
		return nil, fmt.Errorf("failed to create watchdata logs exporter: %w", err)
	}

	return exp, nil
}

// Start is a lifecycle function for the exporter.
func (e *watchdataExporter) Start(ctx context.Context, host component.Host) error {
	_, e.cancel = context.WithCancel(ctx)
	 
	target := fmt.Sprintf("passthrough:///%s", e.endpoint)

	opts := []grpc.DialOption{}
	if e.tlsInsecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		return fmt.Errorf("TLS is required for production")
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return fmt.Errorf("failed to dial gRPC endpoint '%s': %w", e.endpoint, err)
	}

	e.conn = conn
	e.client = plogotlp.NewGRPCClient(conn)

	e.logger.Info("Starting watchdataExporter", zap.String("endpoint", e.endpoint))
	return nil
}

// Shutdown is a lifecycle function for the exporter.
func (e *watchdataExporter) Shutdown(ctx context.Context) error {
	if e.cancel != nil {
		e.cancel()
	}
	if e.conn != nil {
		err := e.conn.Close()
		if err != nil {
			e.logger.Warn("failed to close gRPC connection", zap.Error(err))
		}
	}
	e.logger.Info("Stopping watchdataExporter")
	return nil
}

// Capabilities returns the capabilities of the exporter.
func (e *watchdataExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

// ConsumeLogs is the method that receives log data.
func (e *watchdataExporter) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	return e.sendLogsOverGRPC(ctx, ld)
}

// Compile-time check to ensure watchdataExporter implements exporter.Logs.
// If this line itself causes a compile error, it confirms the interface is not satisfied.
var _ exporter.Logs = (*watchdataExporter)(nil)
