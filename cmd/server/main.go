package main

import (
	"log"
	"net"

	"github.com/Ricky004/watchdata/internals/ingest"
	collectorpb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	listener, err := net.Listen("tcp", ":14317")
	if err != nil {
		log.Fatalf("failed to listen on port 14317: %v", err)
	}

	grpcServer := grpc.NewServer()
	collectorpb.RegisterLogsServiceServer(grpcServer, &ingest.GRPCLogServer{})

	reflection.Register(grpcServer)

	log.Println("WatchData OTLP gRPC server running at :14317")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
