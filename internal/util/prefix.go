package util

import "strings"

// RemovePrefix returns text without a leading color prefix.
func RemovePrefix(text string) string {
	if strings.HasPrefix(text, "[") {
		if idx := strings.Index(text, "]"); idx > 0 {
			return text[idx+1:]
		}
	}
	return text
}

// ApplyPrefix adds or replaces the color prefix.
func ApplyPrefix(text, color string) string {
	base := RemovePrefix(text)
	if color == "" || color == "default" {
		return base
	}
	return "[" + color + "]" + base
}

// ParsePrefix splits a title into color prefix and the remaining text.
func ParsePrefix(text string) (string, string) {
	if strings.HasPrefix(text, "[") {
		if idx := strings.Index(text, "]"); idx > 0 {
			return strings.ToLower(text[1:idx]), text[idx+1:]
		}
	}
	return "", text
}
