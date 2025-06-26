package exporter

func BuildExporter(selected []string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, name := range selected {
		if tmpl, ok := ExporterTemplates[name]; ok {
			result[name] = tmpl
		}
	}
	return result
}