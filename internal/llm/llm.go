package llm

type Message struct {
	User   string
	System string
	Err    error
}
type LLM interface {
	Chat(message Message) (string, error)
}
