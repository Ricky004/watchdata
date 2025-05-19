FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY configs/builder-config.yaml ./configs/builder-config.yaml 

RUN go install go.opentelemetry.io/collector/cmd/builder@latest

RUN builder --config ./configs/builder-config.yaml

FROM alpine:latest
COPY --from=builder /app/dist /dist
COPY configs/otel-collector-config.yaml ./otel-collector-config.yaml
CMD ["/dist/watchdata-collector", "--config", "/otel-collector-config.yaml"] 