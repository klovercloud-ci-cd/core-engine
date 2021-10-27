package v1
type ProcessEvent struct {
	ProcessId string  `bson:"process_id"`
	Data map[string]interface{}  `bson:"data"`
}
