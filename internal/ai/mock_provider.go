package ai

import "context"

type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) Generate(ctx context.Context, input GenerateInput) (*GenerateResult, error) {
	return &GenerateResult{
		Text: `{
  "summary": "AI mock analysis completed for the selected audit logs.",
  "timeline": [
    "Reviewed the filtered audit log set.",
    "Grouped the events into a short incident timeline."
  ],
  "suspicious_signals": [
    "Repeated or clustered actions should be reviewed manually.",
    "Multiple actors or IPs in a narrow time window may indicate unusual activity."
  ],
  "recommendations": [
    "Review the raw audit logs together with user session history.",
    "Confirm whether the affected user or admin activity was expected."
  ]
}`,
	}, nil
}
