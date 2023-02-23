package ansible

type PlayBookResults struct {
	user string
	password string
	apiServer string
	additionalFields map[string]string

}

type PlaybookRunner interface {
	Run()  PlayBookResults
}