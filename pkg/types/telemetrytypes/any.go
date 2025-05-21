package telemetry

type Any struct {
	StringValue *string     `json:"string_value,omitempty"`
	BoolValue   *bool       `json:"bool_value,omitempty"`
	IntValue    *int64      `json:"int_value,omitempty"`
	DoubleValue *float64    `json:"double_value,omitempty"`
	BytesValue  *[]byte     `json:"bytes_value,omitempty"`
	ArrayValue  *[]Any      `json:"array_value,omitempty"`
	KvListValue *[]KeyValue `json:"kv_list_value,omitempty"`
}
