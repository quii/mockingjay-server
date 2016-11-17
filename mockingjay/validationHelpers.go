package mockingjay

import "strings"

func httpHeadersValid(headers map[string]string) bool {
	for k, v := range headers {
		if containsSpace(k) || strings.TrimSpace(v) == "" {
			return false
		}
	}

	return true
}

func containsSpace(s string) bool {
	return strings.Contains(s, " ")
}
