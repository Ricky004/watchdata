package errors


type Error struct {
	Code     Code           `json:"code"`
	Message  string         `json:"message"`
	Severity Severity       `json:"severity"`
	Cause    error          `json:"-"`
	Meta     map[string]any `json:"meta,omitempty"`
}
