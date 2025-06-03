package clickhousestore

import (
	"context"
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/Ricky004/watchdata/pkg/factory"
)

type ClickHouseProvider struct {
	conn clickhouse.Conn
}

func NewClickHouseProvider(ctx context.Context, cfg Config) (*ClickHouseProvider, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
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

	return &ClickHouseProvider{conn: conn}, nil
}

func NewProviderFactory() factory.ProviderFactory[*ClickHouseProvider, Config] {
	return factory.NewProviderFactory(
		factory.MustNewId("clickhouse"),
		func(ctx context.Context, cfg Config) (*ClickHouseProvider, error) {
			return NewClickHouseProvider(ctx, cfg)
		},
	)
}
