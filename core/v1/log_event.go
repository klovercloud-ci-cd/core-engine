package v1

import "time"

type LogEvent struct {
	Index int32 `bson:"index"`
	BuildId      string `bson:"build_id"`
	Log        string   `bson:"log"`
	Step string `bson:"step"`
	CreatedAt    time.Time `bson:"created_at"`
}

func (e LogEvent) Println(s string, s2 string) {

}