package util

import (
	"fmt"
	"os"
	"strings"
)

func GetOpenAIKey() (string, error) {
	key, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		return "", fmt.Errorf("OPENAI_API_KEY is not set")
	}
	return key, nil
}

func ParseArgs(args []string) map[string]string {
	result := make(map[string]string)

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg[2:], "=", 2)
			if len(parts) == 2 {
				result[parts[0]] = parts[1]
			} else {
				key := parts[0]
				if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
					result[key] = args[i+1]
					i++
				} else {
					result[key] = ""
				}
			}
		}
	}

	return result
}
