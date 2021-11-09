package v1

// ProcessEvent Pipeline ProcessEvent struct
type ProcessEvent struct {
	ProcessId string                 `bson:"process_id"`
	Data      map[string]interface{} `bson:"data"`
}
