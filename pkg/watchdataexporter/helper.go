package watchdataexporter

import (
	"github.com/Ricky004/watchdata/pkg/types/telemetrytypes"
	"go.opentelemetry.io/collector/pdata/plog"
)

func convertToLogRecords(ld plog.Logs) []telemetrytypes.LogRecord {
	var records []telemetrytypes.LogRecord
	for i := range ld.ResourceLogs().Len() {
		resLogs := ld.ResourceLogs().At(i)
		for j := range resLogs.ScopeLogs().Len() {
			scopeLogs := resLogs.ScopeLogs().At(j)
			for k := range scopeLogs.LogRecords().Len() {
				log := scopeLogs.LogRecords().At(k)
				records = append(records, telemetrytypes.LogRecord{
					Timestamp:       log.Timestamp().AsTime(),
					ObservedTime:    log.ObservedTimestamp().AsTime(),
					SeverityNumber: int(log.SeverityNumber()),
					SeverityText:   log.SeverityText(),
					Body:            log.Body().AsString(),
					// Add attributes/resource conversion if needed
				})
			}
		}
	}
	return records
}
