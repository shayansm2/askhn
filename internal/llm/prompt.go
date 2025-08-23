package llm

import (
	"fmt"
	"strings"

	"github.com/shayansm2/temporallm/internal/elasticsearch"
)

type SystemPromptBuilder struct{}

func (b SystemPromptBuilder) ForRAG(searchResults []elasticsearch.ESDocument) string {
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

func (b SystemPromptBuilder) ForGrandTruth(path string) string {
	return ""
}
