package clickhousestore

import (
	"fmt"
	"os"
	"time"

	"github.com/Ricky004/watchdata/pkg/factory"
)

type Config struct {
	// Provider is the provider to use
	Provider string `mapstructure:"provider"`

	// Connection is the connection configuration
	Connection ConnectionConfig `mapstructure:",squash"`

	// Clickhouse is the clickhouse configuration
	Clickhouse ClickhouseConfig `mapstructure:"clickhouse"`
}

type ConnectionConfig struct {
	// MaxOpenConns is the maximum number of open connections to the database.
	MaxOpenConns int `mapstructure:"max_open_conns"`

	// MaxIdleConns is the maximum number of connections in the idle connection pool.
	MaxIdleConns int `mapstructure:"max_idle_conns"`

	// DialTimeout is the timeout for dialing a new connection.
	DialTimeout time.Duration `mapstructure:"dial_timeout"`
}

type QuerySettings struct {
	MaxExecutionTime                    int `mapstructure:"max_execution_time"`
	MaxExecutionTimeLeaf                int `mapstructure:"max_execution_time_leaf"`
	TimeoutBeforeCheckingExecutionSpeed int `mapstructure:"timeout_before_checking_execution_speed"`
	MaxBytesToRead                      int `mapstructure:"max_bytes_to_read"`
	MaxResultRowsForCHQuery             int `mapstructure:"max_result_rows_for_ch_query"`
}

type ClickhouseConfig struct {
	// DSN is the database source name.
	DSN string `mapstructure:"dsn"`

	// QuerySettings is the query settings for clickhouse.
	QuerySettings QuerySettings `mapstructure:"settings"`
}

func NewConfigFactory() factory.Factory {
	return factory.NewFactory(factory.MustNewId("clickhousestore"), newConfig)
}

func newConfig() factory.Configurable {
	// Get ClickHouse host from environment or use Docker service name
	clickhouseHost := "clickhouse"
	if host := os.Getenv("CLICKHOUSE_HOST"); host != "" {
		clickhouseHost = host
	}
	
	// Build DSN with proper host
	dsn := fmt.Sprintf("tcp://%s:9000/default?username=default&password=pass", clickhouseHost)
	
	return Config{
		Provider: "clickhouse",
		Connection: ConnectionConfig{
			MaxOpenConns: 100,
			MaxIdleConns: 50,
			DialTimeout:  5 * time.Second,
		},
		Clickhouse: ClickhouseConfig{
			DSN: dsn,
		},
	}
}

func (c Config) Validate() error {
	return nil
}