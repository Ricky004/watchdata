package watchdataexporter

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/zap"
)

type watchdataExporter struct {
	endpoint string
	apiKey string
	logger *zap.Logger
	httpClient *http.Client
}

func newLogsExporter(cfg *Config, set exporter.Settings) (*watchdataExporter, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("endpoint must be provided for watchdataExporter")
	}

	return &watchdataExporter{
		endpoint: cfg.Endpoint,
		apiKey: cfg.APIKey,
		logger: set.Logger,
		httpClient: &http.Client{},
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
	e.logger.Info("Starting watchdataExporter", zap.String("endpoint", e.endpoint))
	return nil
}

// Shutdown is a lifecycle function for the exporter.
func (e *watchdataExporter) Shutdown(ctx context.Context) error {
	e.logger.Info("Stopping watchdataExporter")
	return nil
}

// Capabilities returns the capabilities of the exporter.
func (e *watchdataExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

// ConsumeLogs is the method that receives log data.
func (e *watchdataExporter) ConsumeLogs(ctx context.Context, ld plog.Logs) error {
	return e.sendLogsOverHTTP(ctx, "/logs", ld)
}

// sendLogsOverHTTP is an internal helper method to send data.
func (e *watchdataExporter) sendLogsOverHTTP(ctx context.Context, pathSuffix string, logsData plog.Logs) error {
	// Using OTLP JSON Marshaler for standard OTLP/JSON format
	marshaler := plog.JSONMarshaler{}
	body, err := marshaler.MarshalLogs(logsData)
	if err != nil {
		e.logger.Error("Failed to marshal log data (OTLP/JSON)", zap.Error(err))
		return consumererror.NewPermanent(fmt.Errorf("failed to marshal OTLP logs to JSON: %w", err))
	}

	url := e.endpoint + pathSuffix
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(body))
	if err != nil {
		e.logger.Error("Failed to create HTTP request", zap.String("url", url), zap.Error(err))
		// This error is likely not data-dependent, but could be transient if context is canceled.
		// However, a malformed URL (from endpoint) would be permanent.
		return fmt.Errorf("failed to create HTTP request for %s: %w", url, err)
	}

	req.Header.Set("Content-Type", "application/json")
	if e.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+e.apiKey)
	}
	req.Header.Set("User-Agent", "watchdata-otel-collector-exporter/0.1.0") // Example User-Agent

	e.logger.Debug("Sending log data", zap.String("url", url), zap.Int("payload_size_bytes", len(body)))

	resp, err := e.httpClient.Do(req)
	if err != nil {
		e.logger.Error("Failed to send HTTP request", zap.String("url", url), zap.Error(err))
		return fmt.Errorf("HTTP request to %s failed: %w", url, err) // Retryable for network issues
	}
	defer resp.Body.Close()

	responseBodyBytes, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		e.logger.Warn("Failed to read response body", zap.Int("status_code", resp.StatusCode), zap.Error(readErr))
		// Don't shadow the original error if status code indicates one.
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		e.logger.Debug("Log data sent successfully", zap.Int("status_code", resp.StatusCode), zap.String("url", url))
		return nil // Success
	}

	e.logger.Error("HTTP request to send logs failed",
		zap.Int("status_code", resp.StatusCode),
		zap.String("status", resp.Status),
		zap.String("url", url),
		zap.ByteString("response_body", responseBodyBytes))

	// Classify errors for the collector's retry mechanism
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return consumererror.NewPermanent(fmt.Errorf("permanent auth error (status %d) from %s: %s", resp.StatusCode, url, string(responseBodyBytes)))
	}
	if resp.StatusCode == http.StatusBadRequest || (resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != http.StatusTooManyRequests && resp.StatusCode != http.StatusServiceUnavailable) {
		return consumererror.NewPermanent(fmt.Errorf("permanent client error (status %d) from %s: %s", resp.StatusCode, url, string(responseBodyBytes)))
	}
	// For 5xx server errors, or specific retryable 4xx (like 429 Too Many Requests), return a plain error to signal retry.
	return fmt.Errorf("retryable server/client error (status %d) from %s: %s", resp.StatusCode, url, string(responseBodyBytes))
}

// Compile-time check to ensure watchdataExporter implements exporter.Logs.
// If this line itself causes a compile error, it confirms the interface is not satisfied.
var _ exporter.Logs = (*watchdataExporter)(nil)