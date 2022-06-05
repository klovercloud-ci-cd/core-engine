package v1

// ProcessEvent Pipeline ProcessEvent struct
type ProcessEvent struct {
	ProcessId string                 `bson:"process_id" json:"process_id"`
	CompanyId string  `bson:"company_id" json:"company_id"`
	Data      map[string]interface{} `bson:"data" json:"data"`
}
