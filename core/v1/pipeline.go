package v1

var _ interface {

}= NewPipeline()


type Pipeline struct {
	ApiVersion string `json:"api_version"`
	Name string `json:"name"`
	Label map[string]string  `json:"label"`
	Steps [] Step  `json:"steps"`
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}
