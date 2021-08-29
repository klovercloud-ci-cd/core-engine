package mongo

import (
	"context"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"time"
)

var (
	LogEventCollection="logEventCollection"
)

type logEventRepository struct {
	manager  *dmManager
	timeout  time.Duration
}

func (l logEventRepository) Store(event v1.LogEvent) {
	coll := l.manager.Db.Collection(LogEventCollection)
	_, err := coll.InsertOne(l.manager.Ctx, event)
	if err != nil {
		log.Println("[ERROR] Insert document:", err.Error())
	}
}

func (l logEventRepository) GetByProcessId(processId string, option v1.LogEventQueryOption) []string {
	var results []string
	query:=bson.M{
		"$and": []bson.M{},
	}
	and:=[]bson.M{}
	and= append(and, map[string]interface{}{"process_id": processId})
	if option.IndexTo > 0 {
		and= append(and, map[string]interface{}{"index": bson.M{
			"$gte": option.IndexFrom,
			"$lte": option.IndexTo,
		}})
	}
	if option.Step != "" {
		and= append(and, map[string]interface{}{"step": option.Step})
	}
	query["$and"]=and
	coll := l.manager.Db.Collection(LogEventCollection)
	result, err := coll.Find(l.manager.Ctx, query)
	if err!=nil{
		log.Println(err.Error())
	}
	for result.Next(context.TODO()) {
		elemValue := new(v1.LogEvent)
		err := result.Decode(elemValue)
		if err != nil {
			log.Println("[ERROR]", err)
			break
		}
		results= append(results, elemValue.Log)
	}
	return results
}

func NewLogEventRepository() repository.LogEventRepository {
	return &logEventRepository{

	}

}
