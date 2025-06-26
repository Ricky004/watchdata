package processors

var ProcessorTemplates = map[string]interface{}{
	"batch": map[string]interface{}{
		"send_batch_size": 10000,
		"timeout": "10s",
	},
}