package errors

import (
	"fmt"
	"regexp"
)

type Code string

const (
	// --- Database / Storage Errors ---
	CodeDBConnection      Code = "db.connection_failed"
	CodeDBQueryFailed     Code = "db.query_failed"
	CodeDBMigrationFailed Code = "db.migration_failed"
	CodeCacheUnavailable  Code = "cache.unavailable"
	CodeDBTimeout         Code = "db.timeout"

	// --- Request / Validation Errors ---
	CodeInvalidRequest   Code = "request.invalid"
	CodeMissingField     Code = "request.missing_field"
	CodeInvalidField     Code = "request.invalid_field"
	CodePayloadTooLarge  Code = "request.payload_too_large"
	CodeUnsupportedMedia Code = "request.unsupported_media"

	// --- Authentication / Authorization ---
	CodeUnauthorized Code = "auth.unauthorized"
	CodeForbidden    Code = "auth.forbidden"
	CodeTokenExpired Code = "auth.token_expired"
	CodeTokenInvalid Code = "auth.token_invalid"
	CodeUserDisabled Code = "auth.user_disabled"

	// --- Resource / Object Errors ---
	CodeNotFound          Code = "resource.not_found"
	CodeAlreadyExists     Code = "resource.already_exists"
	CodeConflict          Code = "resource.conflict"
	CodeResourceLocked    Code = "resource.locked"
	CodeDependencyMissing Code = "resource.dependency_missing"

	// --- Rate Limiting / Throttling ---
	CodeRateLimit       Code = "rate_limited"
	CodeTooManyRequests Code = "rate.too_many_requests"

	// --- Internal / Unexpected Errors ---
	CodeInternalError  Code = "internal.error"
	CodePanicRecovered Code = "internal.panic_recovered"
	CodeEncodingFailed Code = "internal.encoding_failed"
	CodeUnknown        Code = "internal.unknown"

	// --- External / 3rd Party Services ---
	CodeExternalAPIError    Code = "external.api_error"
	CodeExternalTimeout     Code = "external.timeout"
	CodeExternalUnreachable Code = "external.unreachable"
	CodeExternalBadResponse Code = "external.bad_response"

	// --- Observability / Traceability ---
	CodeTraceMissing     Code = "observability.trace_missing"
	CodeMetricsCorrupted Code = "observability.metrics_corrupted"
	CodeLogFormatInvalid Code = "observability.log_format_invalid"

	// --- Feature Flags / Config ---
	CodeFeatureDisabled Code = "feature.disabled"
	CodeInvalidConfig   Code = "config.invalid"
	CodeConfigMissing   Code = "config.missing"
)

var validCodePattern = regexp.MustCompile(`^[a-z]+\.[a-z_]+$`)

func validateCodes() error {
	codes := []Code{
		CodeDBConnection,
		CodeDBQueryFailed,
		CodeDBMigrationFailed,
		CodeCacheUnavailable,
		CodeDBTimeout,
		CodeInvalidRequest,
		CodeMissingField,
		CodeInvalidField,
		CodePayloadTooLarge,
		CodeUnsupportedMedia,
		CodeUnauthorized,
		CodeForbidden,
		CodeTokenExpired,
		CodeTokenInvalid,
		CodeUserDisabled,
		CodeNotFound,
		CodeAlreadyExists,
		CodeConflict,
		CodeResourceLocked,
		CodeDependencyMissing,
		CodeTooManyRequests,
		CodeInternalError,
		CodePanicRecovered,
		CodeEncodingFailed,
		CodeUnknown,
		CodeExternalAPIError,
		CodeExternalTimeout,
		CodeExternalUnreachable,
		CodeExternalBadResponse,
		CodeTraceMissing,
		CodeMetricsCorrupted,
		CodeLogFormatInvalid,
		CodeFeatureDisabled,
		CodeInvalidConfig,
		CodeConfigMissing,
	}

	for _, code := range codes {
		if !validCodePattern.MatchString(string(code)) {
			return fmt.Errorf("invalid error code format: %s", code)
		}
	}

	return nil
}

func MustvalidateCodes() {
	if err := validateCodes(); err != nil {
		panic(err)
	}
}
