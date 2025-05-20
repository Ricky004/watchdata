package watchdataexporter

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/pdata/plog"
	plogotlp"go.opentelemetry.io/collector/pdata/plog/plogotlp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (e *watchdataExporter) sendLogsOverGRPC(ctx context.Context, logsData plog.Logs) error {
	conn, err := grpc.NewClient(
	   e.endpoint,
       grpc.WithTransportCredentials(insecure.NewCredentials(),
	), // use TLS in prod
)
	if err != nil {
		return fmt.Errorf("failed to dial gRPC endpoint %s: %w", e.endpoint, err)
	}
	defer conn.Close()

	client := plogotlp.NewGRPCClient(conn)

	// Convert plog.Logs to OTLP ExportLogsServiceRequest
	req := plogotlp.NewExportRequestFromLogs(logsData)

	e.logger.Debug("Sending logs via gRPC", zap.Int("log_record_count", logsData.LogRecordCount()))

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = client.Export(ctx, req)
	if err != nil {
		e.logger.Error("Failed to export logs via gRPC", zap.Error(err))
		return consumererror.NewPermanent(fmt.Errorf("gRPC export failed: %w", err))
	}

	e.logger.Debug("Logs exported via gRPC successfully")
	return nil
}
