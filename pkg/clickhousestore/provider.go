package clickhousestore

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/Ricky004/watchdata/pkg/factory"
	"github.com/Ricky004/watchdata/pkg/types/telemetrytypes"
)

type ClickHouseProvider struct {
	conn clickhouse.Conn
}

func NewClickHouseProvider(ctx context.Context, cfg Config) (*ClickHouseProvider, error) {
	parsedURL, err := url.Parse(cfg.Clickhouse.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ClickHouse DSN: %w", err)
	}
	addr := parsedURL.Host

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
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
		severity_number Int8 CODEC(ZSTD(1)),
		severity_text LowCardinality(String) CODEC(ZSTD(1)),
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
	ORDER BY (timestamp, severity_number)
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

	batch, err := p.conn.PrepareBatch(ctx, "INSERT INTO logs (timestamp, observed_time, severity_number, severity_text, body, attributes, resource, trace_id, span_id, trace_flags, flags, dropped_attributes_count)")
	if err != nil {
		return fmt.Errorf("failed to prepare batch: %w", err)
	}

	for _, log := range logs {
		// Convert attributes and resource to JSON strings
		attributesStr := convertAttributesToString(log.Attributes)
		resourceStr := convertResourceToString(log.Resource)

		err := batch.Append(
			log.Timestamp,
			log.ObservedTime,
			int8(log.SeverityNumber),
			log.SeverityText,
			log.Body,
			attributesStr,
			resourceStr,
			log.TraceID,
			log.SpanID,
			uint8(log.TraceFlags),
			log.Flags,
			uint32(log.DroppedAttrCount),
		)
		if err != nil {
			return fmt.Errorf("failed to append to batch: %w", err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("failed to send batch: %w", err)
	}

	return nil
}

func (p *ClickHouseProvider) GetLogs(ctx context.Context) ([]telemetrytypes.LogRecord, error) {
	var logs []telemetrytypes.LogRecord

	rows, err := p.conn.Query(ctx, `SELECT * FROM logs ORDER BY timestamp DESC LIMIT 1000`)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	defer rows.Close()
	for rows.Next() {
		var log telemetrytypes.LogRecord
		var attributesStr, resourceStr string

		// First scan the raw data from database
		if err := rows.Scan(
			&log.Timestamp,
			&log.ObservedTime,
			&log.SeverityNumber,
			&log.SeverityText,
			&log.Body,
			&attributesStr, // Scan into string variable
			&resourceStr,   // Scan into string variable
			&log.TraceID,
			&log.SpanID,
			&log.TraceFlags,
			&log.Flags,
			&log.DroppedAttrCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		log.Attributes = parseAttributes(attributesStr)
		log.Resource.Attributes = parseResource(resourceStr)

		logs = append(logs, log)
	}

	return logs, nil
}

func (p *ClickHouseProvider) GetLogsSince(ctx context.Context, since time.Time) ([]telemetrytypes.LogRecord, error) {
	query := `SELECT * FROM logs WHERE timestamp > ? ORDER BY timestamp ASC`

	rows, err := p.conn.Query(ctx, query, since)
	if err != nil {
		return nil, fmt.Errorf("failed to query for new logs: %w", err)
	}
	defer rows.Close()

	var logs []telemetrytypes.LogRecord
	for rows.Next() {
		var log telemetrytypes.LogRecord
		var attributesStr, resourceStr string

		if err := rows.Scan(
			&log.Timestamp, &log.ObservedTime, &log.SeverityNumber, &log.SeverityText, &log.Body,
			&attributesStr, &resourceStr, &log.TraceID, &log.SpanID,
			&log.TraceFlags, &log.Flags, &log.DroppedAttrCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan new log row: %w", err)
		}

		log.Attributes = parseAttributes(attributesStr)
		log.Resource.Attributes = parseResource(resourceStr)
		logs = append(logs, log)
	}

	return logs, nil
}

func (p *ClickHouseProvider) GetLogsInTimeRanges(ctx context.Context, startTs, endTs int64) ([]telemetrytypes.LogRecord, error) {
	query := `SELECT * FROM logs
			  WHERE timestamp >= toDateTime(?) AND timestamp <= toDateTime(?)
			  ORDER BY timestamp DESC
			  LIMIT 1000
			  `

	rows, err := p.conn.Query(ctx, query, startTs, endTs)
	if err != nil {
		return nil, fmt.Errorf("failed to query for new logs: %w", err)
	}
	defer rows.Close()

	var logs []telemetrytypes.LogRecord
	for rows.Next() {
		var log telemetrytypes.LogRecord
		var attributesStr, resourceStr string

		if err := rows.Scan(
			&log.Timestamp, &log.ObservedTime, &log.SeverityNumber, &log.SeverityText, &log.Body,
			&attributesStr, &resourceStr, &log.TraceID, &log.SpanID,
			&log.TraceFlags, &log.Flags, &log.DroppedAttrCount,
		); err != nil {
			return nil, fmt.Errorf("failed to scan new log row: %w", err)
		}

		log.Attributes = parseAttributes(attributesStr)
		log.Resource.Attributes = parseResource(resourceStr)
		logs = append(logs, log)
	}

	return logs, nil
}

// helper functions
// Convert []KeyValue to JSON string
func convertAttributesToString(attributes []telemetrytypes.KeyValue) string {
	if len(attributes) == 0 {
		return "{}"
	}

	// Convert to map for cleaner JSON
	attrMap := make(map[string]interface{})
	for _, kv := range attributes {
		attrMap[kv.Key] = kv.Value
	}

	bytes, err := json.Marshal(attrMap)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}

// Convert Resource to JSON string
func convertResourceToString(resource telemetrytypes.Resource) string {
	return convertAttributesToString(resource.Attributes)
}

func parseAttributes(attrStr string) []telemetrytypes.KeyValue {
	if attrStr == "" || attrStr == "{}" {
		return nil
	}
	var attrMap map[string]interface{}
	if json.Unmarshal([]byte(attrStr), &attrMap) != nil {
		return nil
	}
	kvs := make([]telemetrytypes.KeyValue, 0, len(attrMap))
	for k, v := range attrMap {
		kvs = append(kvs, telemetrytypes.KeyValue{Key: k, Value: v})
	}
	return kvs
}

func parseResource(resStr string) []telemetrytypes.KeyValue {
	return parseAttributes(resStr)
}
