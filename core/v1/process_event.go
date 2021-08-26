package v1
type PipelineProcessStatus struct {
	ProcessId string  `bson:"process_id"`
	Data map[string]interface{}  `bson:"data"`
}
