package completion_providers

import "context"

type ICompletionProvider interface {
	GetCompletions(ctx context.Context, completionContext ICompletionProviderContext) []string
}
