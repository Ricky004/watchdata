package telemetrytypes

import "go.opentelemetry.io/collector/pdata/pcommon"

type KeyValue struct {
	Key   string        `json:"key"`
	Value pcommon.Value `json:"value"`
}
