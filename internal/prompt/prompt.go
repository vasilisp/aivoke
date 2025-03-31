package prompt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"text/template"

	"github.com/vasilisp/aivoke/internal/data"
	"github.com/vasilisp/aivoke/internal/util"
)

const dirBasename = ".aivoke"

type Config struct {
	Postprocess bool `json:"postprocess,omitempty"`
}

type Prompt struct {
	Content []byte
	Config  *Config
}

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

func readFromFS(fsys fs.FS, basename string, basenameJSON string) ([]byte, *Config, error) {
	var config Config

	content, err := fs.ReadFile(fsys, basename)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to read prompt: %v", err)
	}

	contentConfig, err := fs.ReadFile(fsys, basenameJSON)
	if err == nil {
		if err := json.Unmarshal(contentConfig, &config); err == nil {
			return content, &config, nil
		}
	}

	return content, &Config{}, nil
}

func read(basename string, template bool) ([]byte, *Config, error) {
	basenameJSON := basename + ".json"
	if template {
		basename = basename + ".tmpl"
	}

	promptDir, err := dir()
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to get prompt directory: %v", err)
	}

	var fsys fs.FS
	if _, err := os.Stat(path.Join(promptDir, basename)); err != nil {
		fsys, err = fs.Sub(data.PromptFS, "prompts")
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to create embedded prompt FS: %v", err)
		}
	} else {
		fsys = os.DirFS(promptDir)
	}

	return readFromFS(fsys, basename, basenameJSON)
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

func Build(id string, args map[string]string) (Prompt, error) {
	var emptyPrompt Prompt

	if err := util.ValidateID(id); err != nil {
		return emptyPrompt, fmt.Errorf("Invalid id: %v", err)
	}

	// prioritize plain prompts over templates and local over embedded

	exists, err := exists(id)
	if err != nil {
		return emptyPrompt, fmt.Errorf("Failed to check if prompt exists: %v", err)
	}

	if exists {
		content, config, err := read(id, false)
		if err != nil {
			return emptyPrompt, fmt.Errorf("Failed to read prompt: %v", err)
		}

		return Prompt{Content: content, Config: config}, nil
	}

	template, config, err := read(id, true)
	if err != nil {
		return emptyPrompt, fmt.Errorf("Failed to read prompt template: %v", err)
	}

	content, err := ExecuteTemplate(template, args)
	if err != nil {
		return emptyPrompt, fmt.Errorf("Failed to execute prompt template: %v", err)
	}

	return Prompt{Content: content, Config: config}, nil
}
