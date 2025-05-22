package errors

type Severity string

const (
	SeverityInfo     Severity = "info"     // harmless events, e.g. optional config fallback
	SeverityWarning  Severity = "warning"  // recoverable issues, e.g. degraded service
	SeverityError    Severity = "error"    // user-impacting or failed operations
	SeverityCritical Severity = "critical" // system-wide issues needing urgent attention
)
