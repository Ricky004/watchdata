
package clickhousestore

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/Ricky004/watchdata/pkg/factory"
	"github.com/Ricky004/watchdata/pkg/types/telemetrytypes"
)

type ClickHouseProvider struct {
	conn clickhouse.Conn
}

func NewClickHouseProvider(ctx context.Context, cfg Config) (*ClickHouseProvider, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"clickhouse:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "pass",
		},
		DialTimeout: cfg.Connection.DialTimeout,
		Settings: map[string]interface{}{
			"max_execution_time":                      cfg.Clickhouse.QuerySettings.MaxExecutionTime,
			"max_execution_time_leaf":                 cfg.Clickhouse.QuerySettings.MaxExecutionTimeLeaf,
			"timeout_before_checking_execution_speed": cfg.Clickhouse.QuerySettings.TimeoutBeforeCheckingExecutionSpeed,
			"max_bytes_to_read":                       cfg.Clickhouse.QuerySettings.MaxBytesToRead,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to open clickhouse connection: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("clickhouse ping failed: %w", err)
	}

	provider := &ClickHouseProvider{conn: conn}

	// Create the logs table if it doesn't exist
	if err := provider.createLogsTable(ctx); err != nil {
		return nil, fmt.Errorf("failed to create logs table: %w", err)
	}

	return provider, nil
}

func (p *ClickHouseProvider) createLogsTable(ctx context.Context) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS logs (
		timestamp DateTime64(9) CODEC(Delta(8), ZSTD(1)),
		observed_time DateTime64(9) CODEC(Delta(8), ZSTD(1)),
		serverity_number Int8 CODEC(ZSTD(1)),
		serverity_text LowCardinality(String) CODEC(ZSTD(1)),
		body String CODEC(ZSTD(1)),
		attributes String CODEC(ZSTD(1)),
		resource String CODEC(ZSTD(1)),
		trace_id FixedString(32) CODEC(ZSTD(1)),
		span_id FixedString(16) CODEC(ZSTD(1)),
		trace_flags UInt8 CODEC(ZSTD(1)),
		flags UInt32 CODEC(ZSTD(1)),
		dropped_attributes_count UInt32 CODEC(ZSTD(1))
	) ENGINE = MergeTree()
	PARTITION BY toYYYYMM(timestamp)
	ORDER BY (timestamp, serverity_number)
	TTL toDateTime(timestamp) + INTERVAL 30 DAY
	SETTINGS index_granularity = 8192, compress_marks = false, compress_primary_key = false;
	`

	if err := p.conn.Exec(ctx, createTableSQL); err != nil {
		return fmt.Errorf("failed to create logs table: %w", err)
	}

	return nil
}

func NewProviderFactory() factory.ProviderFactory[*ClickHouseProvider, Config] {
	return factory.NewProviderFactory(
		factory.MustNewId("clickhouse"),
		func(ctx context.Context, cfg Config) (*ClickHouseProvider, error) {
			return NewClickHouseProvider(ctx, cfg)
		},
	)
}

func (p *ClickHouseProvider) InsertLogs(ctx context.Context, logs []telemetrytypes.LogRecord) error {
	if len(logs) == 0 {
		return nil // Nothing to insert
	}

	batch, err := p.conn.PrepareBatch(ctx, "INSERT INTO logs (timestamp, observed_time, serverity_number, serverity_text, body, attributes, resource, trace_id, span_id, trace_flags, flags, dropped_attributes_count)")
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	for _, log := range logs {
		// Convert attributes and resource to JSON strings if they're not already strings
		attributesStr := convertToString(log.Attributes)
		resourceStr := convertToString(log.Resource)

		err := batch.Append(
			log.Timestamp,
			log.ObservedTime,
			log.ServerityNumber,
			log.ServerityText,
			log.Body,
			attributesStr,
			resourceStr,
			log.TraceID,
			log.SpanID,
			log.TraceFlags,
			log.Flags,
			log.DroppedAttrCount,
		)
		if err != nil {
			return fmt.Errorf("failed to append to batch: %w", err)
		}
	}

	return batch.Send()
}

// Helper function to convert various types to string for ClickHouse storage
func convertToString(v interface{}) string {
	if v == nil {
		return ""
	}
	
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	default:
		// For complex types, you might want to use JSON marshaling
		return fmt.Sprintf("%v", val)
	}
}