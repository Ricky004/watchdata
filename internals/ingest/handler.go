package ingest

import (
    "context"
	"log"

	collectorpb "go.opentelemetry.io/proto/otlp/collector/logs/v1"

)

type GRPCLogServer struct {
	collectorpb.UnimplementedLogsServiceServer
}

func (s *GRPCLogServer) Export(ctx context.Context, req *collectorpb.ExportLogsServiceRequest) (*collectorpb.ExportLogsServiceResponse, error) {
    count := 0
    for _, rl := range req.ResourceLogs {
        // Print resource attributes if available
        if rl.Resource != nil && len(rl.Resource.Attributes) > 0 {
            log.Printf("Resource attributes: %v", rl.Resource.Attributes)
        }
        
        for _, sl := range rl.ScopeLogs {
            // Print scope name if available
            if sl.Scope != nil {
                log.Printf("Scope: %s", sl.Scope.Name)
            }
            
            for i, logRecord := range sl.LogRecords {
                count++
                // Print the actual log content
                log.Printf("Log #%d - Timestamp: %v, Body: %v", 
                    i, 
                    logRecord.TimeUnixNano, 
                    logRecord.Body)
                
                // Print attributes if available
                if len(logRecord.Attributes) > 0 {
                    log.Printf("  Attributes: %v", logRecord.Attributes)
                }
                
                // Print trace context if available
                if logRecord.TraceId != nil || logRecord.SpanId != nil {
                    log.Printf("  Trace ID: %x, Span ID: %x", 
                        logRecord.TraceId, 
                        logRecord.SpanId)
                }
            }
        }
    }
    
    log.Printf("Received %d OTLP log records via gRPC\n", count)
    return &collectorpb.ExportLogsServiceResponse{}, nil
}
