# WatchData Architecture

WatchData is an open-source observability platform built on OpenTelemetry standards, designed to provide lightweight and cost-effective monitoring for logs, metrics, and traces. This document outlines the system architecture, components, and data flow.

## System Overview

WatchData follows a modular architecture with clear separation of concerns:

- **Data Ingestion**: OpenTelemetry Collector with custom exporters
- **Storage**: ClickHouse for high-performance time-series data
- **API Layer**: REST API with WebSocket support for real-time updates
- **Frontend**: Next.js-based web interface for visualization
- **Client Libraries**: gRPC clients for data submission

## Core Components

### 1. OpenTelemetry Collector
The collector serves as the primary data ingestion point, supporting multiple protocols and formats.

**Configuration**: `configs/otel-collector-config.yaml`
- **Receivers**: 
  - OTLP gRPC endpoint (`:4317`) for standard OpenTelemetry data
  - File log receiver for local file ingestion
- **Processors**: Batch processing for efficient data handling
- **Exporters**: Custom WatchData exporter to ClickHouse

### 2. ClickHouse Storage Layer
High-performance columnar database optimized for time-series data.

**Key Features**:
- Optimized schema with compression (ZSTD, Delta encoding)
- Partitioning by month for efficient queries
- TTL-based data retention (30 days default)
- MergeTree engine for fast inserts and queries

**Schema Design**:
```sql
CREATE TABLE logs (
    timestamp DateTime64(9),
    observed_time DateTime64(9),
    severity_number Int8,
    severity_text LowCardinality(String),
    body String,
    attributes String,  -- JSON-encoded key-value pairs
    resource String,    -- JSON-encoded resource attributes
    trace_id FixedString(32),
    span_id FixedString(16),
    trace_flags UInt8,
    flags UInt32,
    dropped_attributes_count UInt32
) ENGINE = MergeTree()
PARTITION BY toYYYYMM(timestamp)
ORDER BY (timestamp, severity_number)
```

### 3. API Server
REST API server providing data access and real-time capabilities.

**Endpoints**:
- `GET /v1/logs` - Retrieve recent logs
- `GET /v1/logs/since` - Get logs since timestamp
- `GET /v1/logs/timerange` - Query logs within time range
- `WebSocket /ws` - Real-time log streaming

**Implementation**: `cmd/server/main.go`
- HTTP server on port `:8080`
- WebSocket support for live updates
- ClickHouse integration via provider pattern

### 4. Frontend Application
Next.js-based web interface for data visualization and exploration.

**Location**: `frontend/`
- React-based components for log visualization
- Real-time updates via WebSocket connection
- TypeScript for type safety
- Tailwind CSS for styling

### 5. Client SDK
gRPC client for programmatic data submission.

**Implementation**: `cmd/client/main.go`
- OpenTelemetry Protocol (OTLP) support
- Direct gRPC communication with collector
- Example usage for testing and integration

## Data Flow

```
[Applications] 
    ↓ (OTLP/gRPC)
[OpenTelemetry Collector]
    ↓ (Custom Exporter)
[ClickHouse Database]
    ↓ (Query API)
[API Server]
    ↓ (REST/WebSocket)
[Frontend Dashboard]
```

### Detailed Flow:

1. **Data Ingestion**:
   - Applications send telemetry data via OTLP gRPC (port 4317)
   - Collector receives and processes data through configured pipeline
   - Custom WatchData exporter transforms and stores data in ClickHouse

2. **Data Storage**:
   - ClickHouse stores logs with optimized schema
   - Data partitioned by month for efficient querying
   - Automatic compression and TTL management

3. **Data Access**:
   - API server queries ClickHouse for log retrieval
   - REST endpoints provide various query patterns
   - WebSocket enables real-time log streaming

4. **Visualization**:
   - Frontend connects to API server
   - Real-time updates via WebSocket connection
   - Interactive log exploration and filtering

## Configuration Management

### Environment Configuration
- Docker Compose orchestration (`docker-compose.yaml`)
- Environment-specific settings via `.env`
- ClickHouse configuration in `clickhouse_config/`

### Application Configuration
- Collector config: `configs/otel-collector-config.yaml`
- ClickHouse users: `configs/clickhouse-users.xml`
- Builder config: `configs/builder-config.yaml`

## Deployment Architecture

### Development Setup
```
┌─────────────────┐    ┌─────────────────┐
│   ClickHouse    │    │   Collector     │
│   (Port 9000)   │◄───│   (Port 4317)   │
│   (Port 8123)   │    │                 │
└─────────────────┘    └─────────────────┘
         ▲                       ▲
         │                       │
┌─────────────────┐    ┌─────────────────┐
│   API Server    │    │   Frontend      │
│   (Port 8080)   │    │   (Port 3000)   │
└─────────────────┘    └─────────────────┘
```

### Production Considerations
- Container orchestration with Docker Compose
- Network isolation with custom bridge network
- Health checks for service dependencies
- Volume persistence for ClickHouse data
- IPv4-only configuration for compatibility

## Key Design Patterns

### Provider Pattern
The system uses a provider pattern for database abstraction:
- `ClickHouseProvider` implements storage operations
- Factory pattern for provider instantiation
- Interface-based design for extensibility

### Configuration Management
- URI-based configuration parsing (`pkg/config/uri.go`)
- Structured configuration with validation
- Environment variable support

### Error Handling
- Comprehensive error wrapping with context
- Graceful degradation for non-critical failures
- Structured logging throughout the system

## Performance Optimizations

### ClickHouse Optimizations
- Column compression with ZSTD and Delta encoding
- Partitioning strategy for time-based queries
- Optimized primary key ordering
- Batch inserts for high throughput

### Application Optimizations
- Connection pooling for database access
- Batch processing in collector pipeline
- WebSocket for efficient real-time updates
- JSON serialization for flexible attribute storage

## Security Considerations

- Database authentication with username/password
- Network isolation in Docker environment
- Input validation for API endpoints
- Secure WebSocket connections

## Monitoring and Observability

The system is designed to be self-monitoring:
- Health checks for all services
- Structured logging with configurable levels
- Connection monitoring and retry logic
- Performance metrics collection capability

## Extensibility

The architecture supports future enhancements:
- Plugin-based exporter system
- Multiple storage backend support
- Additional telemetry data types (metrics, traces)
- Custom visualization components
- Advanced querying and analytics features