package v1

var _ interface {

}= NewPipeline()


type Pipeline struct {
	ApiVersion string `json:"api_version"`
	Name string `json:"name"`
	Label map[string]string  `json:"label"`
}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}
