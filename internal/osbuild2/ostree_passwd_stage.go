package osbuild2

type OSTreePasswdStageOptions struct {
	// OStree ref to create for the commit
	Ref string `json:"ref"`

	// Set the version of the OS as commit metadata
	OSVersion string `json:"os_version,omitempty"`

	// Commit ID of the parent commit
	Parent string `json:"parent,omitempty"`
}

func (OSTreePasswdStageOptions) isStageOptions() {}

type OSTreePasswdStageInput struct {
	inputCommon
	References OSTreePasswdStageReferences `json:"references"`
}

func (OSTreePasswdStageInput) isStageInput() {}

type OSTreePasswdStageInputs struct {
	Tree *OSTreePasswdStageInput `json:"tree"`
}

func (OSTreePasswdStageInputs) isStageInputs() {}

type OSTreePasswdStageReferences []string

func (OSTreePasswdStageReferences) isReferences() {}

// The OSTreePasswdStage (org.osbuild.ostree.passwd) fetches passwd and group files from
// the commit on input and populates the build filesystem tree with /etc/passwd and
// /etc/group to prevent altering UIDs and GIDs.
func NewOSTreePasswdStage(options *OSTreePasswdStageOptions, inputs *OSTreePasswdStageInputs) *Stage {
	return &Stage{
		Type:    "org.osbuild.ostree.passwd",
		Options: options,
		Inputs:  inputs,
	}
}
