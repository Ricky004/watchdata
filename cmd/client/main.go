package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	collectorpb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	"google.golang.org/protobuf/proto"
)

func main() {
	// Construct log record
	now := time.Now()
	nano := uint64(now.UnixNano())

	logRecord := &logspb.LogRecord{
		TimeUnixNano:   nano,
		SeverityText:   "INFO",
		Body:           &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "Hello from OTLP!"}},
		Attributes:     []*commonpb.KeyValue{{Key: "env", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "dev"}}}},
		TraceId:        []byte("0123456789abcdef"),
		SpanId:         []byte("01234567"),
		SeverityNumber: logspb.SeverityNumber_SEVERITY_NUMBER_INFO,
	}

	// Add to request payload
	req := &collectorpb.ExportLogsServiceRequest{
		ResourceLogs: []*logspb.ResourceLogs{
			{
				Resource: &resourcepb.Resource{
					Attributes: []*commonpb.KeyValue{
						{Key: "service.name", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "watchdata"}}},
					},
				},
				ScopeLogs: []*logspb.ScopeLogs{
					{
						LogRecords: []*logspb.LogRecord{logRecord},
					},
				},
			},
		},
	}

	// Marshal to protobuf
	data, err := proto.Marshal(req)
	if err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}

	// Log the size of the marshaled data.  This can help in debugging.
	log.Printf("Marshaled data size: %d bytes", len(data))

	// Send HTTP request to OpenTelemetry Collector
	reqBody := bytes.NewReader(data) // Create a reader from the byte slice
	resp, err := http.Post("http://localhost:4319/v1/logs", "application/x-protobuf", reqBody)
	if err != nil {
		log.Fatalf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the entire response body.  This is important for error handling.
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err) // Log, but don't fatal
	}

	fmt.Println("Sent log. Collector responded with:", resp.Status)
	fmt.Printf("Response Body: %s\n", respBody) // Print the response body

	if resp.StatusCode != http.StatusAccepted {
		log.Printf("Collector returned non-202 status: %s", resp.Status)
	}
}