package receviers

var ReceiverTemplates = map[string]interface{}{
	"otlp": map[string]interface{}{
		"protocols": map[string]interface{}{
			"grpc": map[string]interface{}{
				"endpoint": "0.0.0.0:4317",
			},
		},
	},
	
	"filelog": map[string]interface{}{
		"include":  []string{"../tmp/test.json"},
		"start_at": "beginning",
		"operators": []map[string]interface{}{
			{
				"type": "json_parser",
				"timestamp": map[string]interface{}{
					"parse_from": "attributes.time",
					"layout":     "%Y-%m-%d %H:%M:%S",
				},
			},
		},
	},
}
