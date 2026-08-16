package completion_providers

type ICompletionProviderContext interface {
	PartialInput() string
	Submodule() interface{}
	Filter() string
}

type CompletionProviderContext struct {
	partialInput string
	submodule    interface{}
	filter       string
}

func NewCompletionProviderContext(partialInput string, submodule interface{}, filter string) *CompletionProviderContext {
	return &CompletionProviderContext{
		partialInput: partialInput,
		submodule:    submodule,
		filter:       filter,
	}
}

func (c *CompletionProviderContext) PartialInput() string {
	return c.partialInput
}

func (c *CompletionProviderContext) Submodule() interface{} {
	return c.submodule
}

func (c *CompletionProviderContext) Filter() string {
	return c.filter
}
