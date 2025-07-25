package main

import (
	"context"
	"encoding/hex"
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
	// Connect to OpenTelemetry gRPC receiver
	conn, err := grpc.NewClient("0.0.0.0:4317", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to gRPC endpoint: %v", err)
	}
	defer conn.Close()

	client := collectorpb.NewLogsServiceClient(conn)

	// Prepare log record data
	now := time.Now()
	nano := uint64(now.UnixNano())

	traceID, _ := hex.DecodeString("2123252589abcdef0123456785abcdef") // 32 hex chars
	spanID, _ := hex.DecodeString("2124466449abcdef")                  // 16 hex chars

	for range 1 {
		logRecord := &logspb.LogRecord{
			TimeUnixNano:         nano,
			ObservedTimeUnixNano: nano,
			SeverityNumber:       logspb.SeverityNumber_SEVERITY_NUMBER_ERROR,
			SeverityText:         "CRITICAL",
			Body: &commonpb.AnyValue{
				Value: &commonpb.AnyValue_StringValue{StringValue: "K8s return a server 505 error to the request you send that not reach there to get the system call"},
			},
			Attributes: []*commonpb.KeyValue{
				{Key: "env", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "dev"}}},
				{Key: "host", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "local"}}},
			},
			TraceId:                traceID,
			SpanId:                 spanID,
			Flags:                  4,
			DroppedAttributesCount: 1,
		}

		// Create resource and scope wrapper
		resourceLogs := &logspb.ResourceLogs{
			Resource: &resourcepb.Resource{
				Attributes: []*commonpb.KeyValue{
					{Key: "source.name", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "nodejs"}}},
					{Key: "service.name", Value: &commonpb.AnyValue{Value: &commonpb.AnyValue_StringValue{StringValue: "frontend"}}},
				},
			},
			ScopeLogs: []*logspb.ScopeLogs{
				{
					LogRecords: []*logspb.LogRecord{logRecord},
				},
			},
		}

		req := &collectorpb.ExportLogsServiceRequest{
			ResourceLogs: []*logspb.ResourceLogs{resourceLogs},
		}

		resp, err := client.Export(context.Background(), req)
		if err != nil {
			log.Fatalf("Failed to export logs via gRPC: %v", err)
		}

		log.Printf("Successfully exported logs: %+v", resp)
	} 

}
