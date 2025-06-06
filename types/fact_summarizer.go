package types

import "context"

type SummarizerFunction func(ctx context.Context, facts map[string]string) (string, error)
