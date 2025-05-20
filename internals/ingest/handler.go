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
		for _, sl := range rl.ScopeLogs {
			count += len(sl.LogRecords)
		}
	}
	log.Printf("Received %d OTLP log records via gRPC\n", count)

	return &collectorpb.ExportLogsServiceResponse{}, nil
}
