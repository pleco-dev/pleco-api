package ai

type GenerateInput struct {
	SystemPrompt string
	UserPrompt   string
	Model        string
	Temperature  float64
	MaxTokens    int
}

type GenerateResult struct {
	Text string
}
