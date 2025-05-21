package utils

import "strings"

// truncateBySeparator truncates the string, keeping up to maxParts segments,
// appending "..." if there are more segments than maxParts.
func truncateBySeparator(s, sep string, maxParts int) string {
	parts := strings.Split(s, sep)
	if len(parts) <= maxParts {
		return s
	}
	return strings.Join(parts[:maxParts], sep) + "..."
}

// TruncateSmart truncates only the content inside parentheses if present,
// otherwise truncates the entire string.
func TruncateSmart(s, sep string, maxParts int) string {
	start := strings.Index(s, "(")
	end := strings.LastIndex(s, ")")
	if start != -1 && end != -1 && start < end {
		prefix := s[:start+1]
		content := s[start+1 : end]
		suffix := s[end:]
		return prefix + truncateBySeparator(content, sep, maxParts) + suffix
	}
	return truncateBySeparator(s, sep, maxParts)
}
