package internal

func IsValidFieldName(field string) bool {
	// Implement validation logic here
	// This is a simple example - adjust based on your schema naming conventions
	for _, char := range field {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') || char == '_' || char == '.') {
			return false
		}
	}
	return len(field) > 0
}
