package prompt

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"text/template"

	"github.com/vasilisp/aivoke/internal/data"
	"github.com/vasilisp/aivoke/internal/util"
)

const dirBasename = ".aivoke"

func dir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("Failed to get home directory: %v", err)
	}

	return path.Join(homeDir, dirBasename), nil
}

func fullLocalPath(basename string) (string, error) {
	promptDir, err := dir()
	if err != nil {
		return "", fmt.Errorf("Failed to get prompt directory: %v", err)
	}

	return path.Join(promptDir, basename), nil
}

func fileExists(fsys fs.FS, path string) bool {
	file, err := fsys.Open(path)
	if err != nil {
		return false
	}

	file.Close()
	return true
}

func exists(basename string) (bool, error) {
	promptPath, err := fullLocalPath(basename)
	if err != nil {
		return false, fmt.Errorf("Failed to get prompt path: %v", err)
	}

	if _, err := os.Stat(promptPath); err == nil {
		return true, nil
	}

	return fileExists(data.PromptFS, path.Join("prompts", basename)), nil
}

func read(basename string, template bool) ([]byte, error) {
	if template {
		basename = basename + ".tmpl"
	}

	promptPath, err := fullLocalPath(basename)
	if err != nil {
		return nil, fmt.Errorf("Failed to get prompt path: %v", err)
	}

	if _, err := os.Stat(promptPath); err != nil {
		return data.PromptFS.ReadFile(path.Join("prompts", basename))
	}

	return os.ReadFile(promptPath)
}

func ExecuteTemplate(tmplBytes []byte, args map[string]string) ([]byte, error) {
	tmpl, err := template.New("prompt").Parse(string(tmplBytes))
	if err != nil {
		return nil, fmt.Errorf("Failed to parse prompt: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, args); err != nil {
		return nil, fmt.Errorf("Failed to execute prompt: %v", err)
	}

	return buf.Bytes(), nil
}

func Build(id string, args map[string]string) ([]byte, error) {
	if err := util.ValidateID(id); err != nil {
		return nil, fmt.Errorf("Invalid id: %v", err)
	}

	// prioritize plain prompts over templates and local over embedded

	exists, err := exists(id)
	if err != nil {
		return nil, fmt.Errorf("Failed to check if prompt exists: %v", err)
	}

	if exists {
		content, err := read(id, false)
		if err != nil {
			return nil, fmt.Errorf("Failed to read prompt: %v", err)
		}

		return content, nil
	}

	template, err := read(id, true)
	if err != nil {
		return nil, fmt.Errorf("Failed to read prompt template: %v", err)
	}

	return ExecuteTemplate(template, args)
}
