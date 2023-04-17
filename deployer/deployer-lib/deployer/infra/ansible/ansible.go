package ansible

type Playbook struct{
	// Name of playbook
	Name string
	//Path to Playbook
	Path string
	OverrideParametersPath string
	RenderedTemplatePath string
	OutputLocation string
}
