package ansible

type PlayBookResults struct {
	results map[string]string

}

type PlaybookRunner interface {
	Run()  PlayBookResults
}