package llm

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Message struct {
	User   string
	System string
	Err    error
}
type LLM interface {
	Chat(message Message) (string, error)
}

func ChatAndCheck(llm LLM, message Message, check func(string) error) (string, error) {
	response, err := llm.Chat(message)
	if err != nil {
		return "", err
	}
	err = check(response)
	if err != nil {
		return "", err
	}
	return response, nil
}

func ChatWithSchema(llm LLM, message Message, v any) error {
	response, err := llm.Chat(message)
	if err != nil {
		return err
	}
	response = strings.TrimSpace(response)
	re := regexp.MustCompile("^```json\\s*|\\s*```$")
	response = re.ReplaceAllString(response, "")

	if err := json.Unmarshal([]byte(response), v); err != nil {
		return errors.Join(err, fmt.Errorf(response))
	}
	return nil
}
