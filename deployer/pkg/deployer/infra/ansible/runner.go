package ansible

type PlayBookResults struct {
	User      string
	Password         string
	ApiServer        string
	AdditionalFields map[string]string

}

type PlaybookRunner interface {
	Run()  PlayBookResults
}