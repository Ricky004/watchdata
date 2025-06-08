package watchdataexporter

import (
	"github.com/Ricky004/watchdata/pkg/types/telemetrytypes"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
)

func convertToLogRecords(ld plog.Logs) []telemetrytypes.LogRecord {
	var records []telemetrytypes.LogRecord
	
	for i := range ld.ResourceLogs().Len() {
		resLogs := ld.ResourceLogs().At(i)
		resource := resLogs.Resource()
		
		// Convert resource attributes
		resourceAttrs := make([]telemetrytypes.KeyValue, 0)
		resource.Attributes().Range(func(k string, v pcommon.Value) bool {
			resourceAttrs = append(resourceAttrs, telemetrytypes.KeyValue{
				Key:   k,
				Value: v.AsString(),
			})
			return true
		})
		
		for j := range resLogs.ScopeLogs().Len() {
			scopeLogs := resLogs.ScopeLogs().At(j)
			
			for k := range scopeLogs.LogRecords().Len() {
				log := scopeLogs.LogRecords().At(k)
				
				// Convert log attributes
				logAttrs := make([]telemetrytypes.KeyValue, 0)
				log.Attributes().Range(func(k string, v pcommon.Value) bool {
					logAttrs = append(logAttrs, telemetrytypes.KeyValue{
						Key:   k,
						Value: v.AsString(),
					})
					return true
				})
				
				// Convert TraceID and SpanID properly
				traceID := ""
				spanID := ""
				if !log.TraceID().IsEmpty() {
					traceID = log.TraceID().String() // This gives hex string
				}
				if !log.SpanID().IsEmpty() {
					spanID = log.SpanID().String() // This gives hex string
				}
				
				records = append(records, telemetrytypes.LogRecord{
					Timestamp:        log.Timestamp().AsTime(),
					ObservedTime:     log.ObservedTimestamp().AsTime(),
					SeverityNumber:   int(log.SeverityNumber()),
					SeverityText:     log.SeverityText(),
					Body:             log.Body().AsString(),
					Attributes:       logAttrs,
					Resource: telemetrytypes.Resource{
						Attributes: resourceAttrs,
					},
					TraceID:          traceID,
					SpanID:           spanID,
					TraceFlags:       uint32(log.Flags()),
					Flags:            uint32(log.Flags()),
					DroppedAttrCount: int(log.DroppedAttributesCount()),
				})
			}
		}
	}
	return records
}


