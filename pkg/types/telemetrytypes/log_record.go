package telemetrytypes

import (
	"time"
)

type Resource struct {
	Attributes []KeyValue `json:"attributes"`
}

type LogRecord struct {
	Timestamp        time.Time  `json:"timestamp"`
	ObservedTime     time.Time  `json:"observed_time"`
	SeverityNumber   int8       `json:"severity_number"`
	SeverityText     string     `json:"severity_text"`
	Body             string     `json:"body"`
	Attributes       []KeyValue `json:"attributes"`
	Resource         Resource   `json:"resource"`
	TraceID          string     `json:"trace_id,omitempty"`
	SpanID           string     `json:"span_id,omitempty"`
	TraceFlags       uint8      `json:"trace_flags,omitempty"`
	Flags            uint32     `json:"flags,omitempty"`
	DroppedAttrCount uint32     `json:"dropped_attributes_count,omitempty"`
}
