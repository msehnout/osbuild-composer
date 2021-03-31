package osbuild2

// Inputs for ostree commits
type OSTreeCheckoutInput struct {
	inputCommon
}

func (OSTreeCheckoutInput) isInput() {}

func NewOSTreeCheckoutInput() *OSTreeCheckoutInput {
	input := new(OSTreeCheckoutInput)
	input.Type = "org.osbuild.ostree.checkout"
	input.Origin = "org.osbuild.source"
	return input
}
