package v1

type Pipeline struct {
	Option PipelineApplyOption
	ApiVersion string `json:"api_version"`
	Name string `json:"name"`
	BuildId string  `json:"build_id"`
	Label map[string]string  `json:"label"`
	Steps [] Step  `json:"steps"`
}
func(pipeline Pipeline)Validate()error{
	for _,each:=range pipeline.Steps{
		err:=each.Validate()
		if err!=nil{
			return err
		}
	}

	return nil
}