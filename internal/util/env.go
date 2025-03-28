package util

import (
	"fmt"
	"os"
)

func GetOpenAIKey() (string, error) {
	key, ok := os.LookupEnv("OPENAI_API_KEY")
	if !ok {
		return "", fmt.Errorf("OPENAI_API_KEY is not set")
	}
	return key, nil
}
