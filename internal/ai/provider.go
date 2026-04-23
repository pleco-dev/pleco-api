package ai

import "context"

type Provider interface {
	Generate(ctx context.Context, input GenerateInput) (*GenerateResult, error)
}
