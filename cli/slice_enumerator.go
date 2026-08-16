package cli

// sliceEnumerator walks a slice of strings similar to C#'s IEnumerator<string>.
type sliceEnumerator struct {
	tokens []string
	index  int
}

func newSliceEnumerator(tokens []string) *sliceEnumerator {
	return &sliceEnumerator{tokens: tokens, index: -1}
}

func (this *sliceEnumerator) MoveNext() bool {
	this.index++
	return this.index < len(this.tokens)
}

func (this *sliceEnumerator) Current() string {
	if this.index < 0 || this.index >= len(this.tokens) {
		return ""
	}
	return this.tokens[this.index]
}
