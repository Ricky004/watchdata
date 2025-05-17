# Variables
BINARY_NAME=watchdata
BACKEND_DIR=cmd/server
COLLECTOR_CONFIG=configs/otel-collector-config.yaml

# Build the Go backend
build:
	go build -o $(BINARY_NAME) $(BACKEND_DIR)

# Run the Go backend
run:
	go run $(BACKEND_DIR)/main.go

# Run the OpenTelemetry Collector
collector:
	docker run --rm -v $(PWD)/$(COLLECTOR_CONFIG):/etc/otel/config.yaml \
	  -p 4319:4319 -p 4318:4318 \
	  otel/opentelemetry-collector:latest \
	  --config=/etc/otel/config.yaml

# Format Go code
fmt:
	go fmt ./...

# Run tests
test:
	go test ./...

# Clean build files
clean:
	rm -f $(BINARY_NAME)

# All
all: fmt build run
