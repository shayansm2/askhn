package llm

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
)

func GenerateSysPrompt(promptName string, values map[string]string) (string, error) {
	pwd, _ := os.Getwd()
	content, err := os.ReadFile(fmt.Sprintf("%v/internal/llm/prompts/%v.txt", pwd, promptName))
	if err != nil {
		return "", err
	}

	promptTemplate, err := template.New(promptName).Parse(string(content))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = promptTemplate.Execute(&buf, values)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
