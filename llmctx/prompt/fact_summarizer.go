package prompt

import (
	"context"
)

//go:generate mockgen -source=fact_summarizer.go -destination=../mocks/mock_fact_summarizer.go -package=mocks
type FactSummarizer interface {
	// Summarizes user facts into a compact string (e.g., for prompt context)
	Summarizer(ctx context.Context, facts map[string]string) (string, error)
}

func DefaultSummarizer(ctx context.Context, facts map[string]string) (string, error) {
	if len(facts) == 0 {
		return "No facts available.", nil
	}

	summary := "User Facts:\n"
	for key, value := range facts {
		summary += "- " + key + ": " + value + "\n"
	}
	return summary, nil
}