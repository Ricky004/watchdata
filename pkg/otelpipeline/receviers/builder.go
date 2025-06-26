package receviers

func BuildReceivers(selected []string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, name := range selected {
		if tmpl, ok := ReceiverTemplates[name]; ok {
			result[name] = tmpl
		}
	}
	return result
}
