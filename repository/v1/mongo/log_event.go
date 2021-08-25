package mongo

import (
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (l logEventRepository) Store(log v1.LogEvent) {
	coll := l.manager.Db.Collection(LogEventCollection)
	_, err := coll.InsertOne(l.manager.Ctx, log)
	if err != nil {
		log.Println("[ERROR] Insert document:", err.Error())
	}
}

func (l logEventRepository) GetByBuildId(buildId string, option v1.LogEventQueryOption) []string {
	query := bson.M{
		"$and": []bson.M{
			{"build_id": buildId},
		},
	}
	findOptions := options.Find()
	if option.IndexTo > 0 {
		findOptions.SetMin(option.IndexTo)
	}
	if option.IndexFrom > 0 {
		findOptions.SetMax(option.IndexFrom)
	}
	var results []string
	coll := l.manager.Db.Collection(LogEventCollection)
	result, _ := coll.Find(l.manager.Ctx, findOptions)
	if option.Step != "" {
		query := bson.M{
			"$and": []bson.M{
				{"build_id": buildId},
				{"step": option.Step},
			},
		}
		coll := l.manager.Db.Collection(LogEventCollection)
		result := coll.FindOne(l.manager.Ctx, query)
		err := result.Decode(results)
		if err != nil {
			log.Println("[ERROR]", err)
		}
		results = append(results)
		return results
	}
	if option.Step == "" && option.IndexTo ==  0 && option.IndexFrom == 0{
		result := coll.FindOne(l.manager.Ctx, query)
		err := result.Decode(results)
		if err != nil {
			log.Println("[ERROR]", err)
		}
		results = append(results)
		return results
	}
	err := result.Decode(results)
	if err != nil {
		log.Println("[ERROR]", err)
	}
	results = append(results)
	return results
}

func NewLogEventRepository() repository.LogEventRepository {
	return &logEventRepository{
	}

}
