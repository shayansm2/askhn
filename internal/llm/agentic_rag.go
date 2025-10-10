package llm

type AgentAction string

const (
	ActionSearch  AgentAction = "SEARCH"
	ActionExplore AgentAction = "EXPLORE"
	ActionAnswer  AgentAction = "ANSWER"
)

type AgentResponse struct {
	Action  AgentAction `json:"action"`
	Message string      `json:"message"`
}
