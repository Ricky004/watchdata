package processors

func BuildProcessors(selected []string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, name := range selected {
		if tmpl, ok := ProcessorTemplates[name]; ok {
			result[name] = tmpl
		}
	}
	return result
}