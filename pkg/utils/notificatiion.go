package utils

func replaceTemplatePlaceholders(message string, details map[string]string) string {
	for key, value := range details {
		placeholder := "{{" + key + "}}"
		message = strings.ReplaceAll(message, placeholder, value)
	}
	return message
}