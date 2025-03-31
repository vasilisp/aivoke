package prompt

import "strings"

func removeCodeFences(response string) string {
	lines := strings.Split(response, "\n")

	start := 0
	end := len(lines) - 1

	// Find first non-empty line
	ignoreSpace := func() {
		for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
			start++
		}
		for end >= 0 && strings.TrimSpace(lines[end]) == "" {
			end--
		}
	}

	ignoreSpace()

	if start < len(lines) &&
		strings.HasPrefix(strings.TrimSpace(lines[start]), "```") &&
		end >= 0 &&
		strings.TrimSpace(lines[end]) == "```" {
		end--
	}

	ignoreSpace()

	if start <= end {
		return strings.Join(lines[start:end+1], "\n")
	}

	return response
}

func Postprocess(response string) string {
	return removeCodeFences(response)
}
