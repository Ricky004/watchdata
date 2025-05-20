package main

import (
	"context"
	"log"
	"time"

	collectorpb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	resourcepb "go.opentelemetry.io/proto/otlp/resource/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// connect to grpc endpoint (4317)
	conn, err := grpc.NewClient("0.0.0.0:4317", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC: %v", err)
	}
	defer conn.Close()

	client := collectorpb.NewLogsServiceClient(conn)

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

	// Send log export
	resp, err := client.Export(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to export logs via gRPC: %v", err)
	}

	log.Printf("Exported logs via gRPC: %+v", resp)
}
