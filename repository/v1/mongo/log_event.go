package mongo

import (
	"context"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// LogEventCollection collection name
var (
	LogEventCollection = "logEventCollection"
)

type logEventRepository struct {
	manager *dmManager
	timeout time.Duration
}

func (l logEventRepository) Store(event v1.LogEvent) {
	coll := l.manager.Db.Collection(LogEventCollection)
	_, err := coll.InsertOne(l.manager.Ctx, event)
	if err != nil {
		log.Println("[ERROR] Insert document:", err.Error())
	}
}

func (l logEventRepository) GetByProcessId(processId string, option v1.LogEventQueryOption) ([]string, int64) {
	var results []string
	query := bson.M{
		"$and": []bson.M{},
	}
	and := []bson.M{{"process_id": processId}}
	if option.Step != "" {
		and = append(and, map[string]interface{}{"step": option.Step})
	}
	query["$and"] = and
	coll := l.manager.Db.Collection(LogEventCollection)
	skip := option.Pagination.Page * option.Pagination.Limit
	result, err := coll.Find(l.manager.Ctx, query, &options.FindOptions{
		Limit: &option.Pagination.Limit,
		Skip:  &skip,
	})
	if err != nil {
		log.Println(err.Error())
	}
	for result.Next(context.TODO()) {
		elemValue := new(v1.LogEvent)
		err := result.Decode(elemValue)
		if err != nil {
			log.Println("[ERROR]", err)
			break
		}
		results = append(results, elemValue.Log)
	}
	count, err := coll.CountDocuments(l.manager.Ctx, query)
	if err != nil {
		log.Println(err.Error())
	}
	return results, count
}

// NewLogEventRepository returns LogEventRepository type object
func NewLogEventRepository(timeout int) repository.LogEventRepository {
	return &logEventRepository{
		manager: GetDmManager(),
		timeout: time.Duration(timeout),
	}

}
