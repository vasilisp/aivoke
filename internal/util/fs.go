package util

import (
	"fmt"
	"os"
	"path"
)

func promptDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Failed to get home directory: %v", err)
	}

	return path.Join(homeDir, ".aivoke"), nil
}

func PromptOfId(id string) ([]byte, error) {
	if err := ValidateID(id); err != nil {
		return nil, fmt.Errorf("Invalid ID: %v", err)
	}

	promptDir, err := promptDir()
	if err != nil {
		return nil, fmt.Errorf("Failed to get prompt directory: %v", err)
	}
	promptPath := path.Join(promptDir, id)

	content, err := os.ReadFile(promptPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to read prompt file: %v", err)
	}

	return content, nil
}
