package ansible

type Playbook struct{
	// Name of playbook
	name string
	//Path to Playbook
	path string
	overrideParametersPath string
	renderedTemplatePath string
}
