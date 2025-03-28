package util

import (
	"os"
)

func GetOpenAIKey() string {
	return os.Getenv("OPENAI_API_KEY")
}
