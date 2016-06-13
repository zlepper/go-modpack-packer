package helpers

import "strings"

func IgnoreCaseContains(s, content string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(content))
}
