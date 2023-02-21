package ansible

type PlayBookResults struct {
	results map[string]string

}

type PlaybookRunner interface {
	run(pathToPlaybook string, pathToParametersFile string)  PlayBookResults
}