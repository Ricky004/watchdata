package errors

import (
	"errors"
	"fmt"
	"log/slog"
)

// Error interface compatibility
func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}

	return fmt.Sprintf("%s[%s]: %s", e.Code, e.Severity, e.Message)
}

// Unwrap for 'errors.Is/As'
func (e *Error) Unwrap() error {
	return e.Cause
}


func (e *Error) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("code", string(e.Code)),
		slog.String("severity", string(e.Severity)),
		slog.String("message", e.Message),
		slog.Any("meta", e.Meta),
	)
}


// New creates a new structured error
func New(code Code, message string, severity Severity, cause ...error) *Error {
	var root error
	if len(cause) > 0 {
		root = cause[0]
	}
	return &Error{
		Code:     code,
		Message:  message,
		Severity: severity,
		Cause: root,
	}
}

// NewMeta includes structured meta data
func NewMeta(code Code, message string, severity Severity, meta map[string]any, cause ...error) *Error {
	var root error
	if len(cause) > 0 {
		root = cause[0]
	}
	return &Error{
		Code:     code,
		Message:  message,
		Severity: severity,
		Cause:    root,
		Meta:     meta,
	}
}

func Wrap(err error, code Code, message string, severity Severity) *Error {
	var e *Error
	if errors.As(err, &e) {
		return &Error{
			Code: code,
			Message: message,
			Severity: severity,
			Cause: e,
			Meta: e.Meta,
		}
	} 
    
	return &Error{
		Code:     code,
		Message:  message,
		Severity: severity,
		Cause:    err,
	}
}
