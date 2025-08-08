package chatbot

import (
	"fmt"
	"strings"
)

func BuildSystemPrompt(searchResults []ElasticSearchDocument) string {
	const promptTemplate = `You're a course teaching assistant. Answer the questions based on the CONTEXT from the hacker news comments.
	Use only the facts from the CONTEXT when answering the questions.

	<CONTEXT>
	%s
	</CONTEXT>`

	var contextBuilder strings.Builder
	for _, doc := range searchResults {
		contextBuilder.WriteString(fmt.Sprintf("title: %s\ncomment: %s\n\n", doc.Title, doc.Text))
	}
	context := strings.TrimSpace(contextBuilder.String())
	prompt := fmt.Sprintf(promptTemplate, context)
	return strings.TrimSpace(prompt)
}
