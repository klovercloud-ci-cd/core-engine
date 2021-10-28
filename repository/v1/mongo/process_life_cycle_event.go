package mongo

import (
	"context"
	v1 "github.com/klovercloud-ci/core/v1"
	"github.com/klovercloud-ci/core/v1/repository"
	"github.com/klovercloud-ci/enums"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	ProcessLifeCycleCollection="processLifeCycleEventCollection"
)


type processLifeCycleRepository struct {
	manager  *dmManager
	timeout  time.Duration
}

func (p processLifeCycleRepository) PullNonInitializedAndAutoTriggerEnabledEventsByStepType(count int64, stepType string) []v1.ProcessLifeCycleEvent {
	var data []v1.ProcessLifeCycleEvent
	query:=bson.M{
		"$and": []bson.M{
			{"status":enums.NON_INITIALIZED},
			{"trigger":enums.AUTO},
			{"step_type":stepType},

		},
	}
	coll := p.manager.Db.Collection(ProcessLifeCycleCollection)
	result, err := coll.Find(p.manager.Ctx, query,&options.FindOptions{
		Limit: &count,
	})
	if err!=nil{
		log.Println(err.Error())
	}
	for result.Next(context.TODO()) {
		elemValue := new(v1.ProcessLifeCycleEvent)
		err := result.Decode(elemValue)
		if err != nil {
			log.Println("[ERROR]", err)
			break
		}
		data= append(data, *elemValue)
	}
	for _,each:=range data{
		go p.updateStatus(each,string(enums.ACTIVE))
	}
	return data
}

func (p processLifeCycleRepository) PullPausedAndAutoTriggerEnabledResourcesByAgentName(count int64,agent string) []v1.ProcessLifeCycleEvent {
	var data []v1.ProcessLifeCycleEvent
	query:=bson.M{
		"$and": []bson.M{
			{"agent":agent},
			{"status":enums.PAUSED},
			{"trigger":enums.AUTO},
			{"step_type":enums.DEPLOY},

		},
	}
	coll := p.manager.Db.Collection(ProcessLifeCycleCollection)
	result, err := coll.Find(p.manager.Ctx, query,&options.FindOptions{
		Limit: &count,
	})
	if err!=nil{
		log.Println(err.Error())
	}
	for result.Next(context.TODO()) {
		elemValue := new(v1.ProcessLifeCycleEvent)
		err := result.Decode(elemValue)
		if err != nil {
			log.Println("[ERROR]", err)
			break
		}
		data= append(data, *elemValue)
	}
	for _,each:=range data{
		go p.updateStatus(each,string(enums.ACTIVE))
	}
	return data
}

func (p processLifeCycleRepository) Get(count int64) []v1.ProcessLifeCycleEvent {
	var data []v1.ProcessLifeCycleEvent
	query:=bson.M{
		"$and": []bson.M{

		},
	}
	coll := p.manager.Db.Collection(ProcessLifeCycleCollection)
	result, err := coll.Find(p.manager.Ctx, query,&options.FindOptions{
		Limit: &count,
	})
	if err!=nil{
		log.Println(err.Error())
	}
	for result.Next(context.TODO()) {
		elemValue := new(v1.ProcessLifeCycleEvent)
		err := result.Decode(elemValue)
		if err != nil {
			log.Println("[ERROR]", err)
			break
		}
		data= append(data, *elemValue)
	}
	return data
}

func (p processLifeCycleRepository) Store(events []v1.ProcessLifeCycleEvent) {
	coll := p.manager.Db.Collection(ProcessLifeCycleCollection)
	var pipeline *v1.Pipeline
	if len(events)>0{
		if events[0].StepType==enums.BUILD{
			pipeline=events[0].Pipeline
		}
	}
	for _,each:= range events{
		existing:=p.GetByProcessIdAndStep(each.ProcessId,each.Step)
		if existing==nil{
			each.Pipeline=pipeline
			_,err:=coll.InsertOne(p.manager.Ctx,each)
			if err!=nil{
				log.Println(err.Error())
			}
		}else{
			existing.Status=each.Status
			err:=p.update(*existing)
			if err!=nil{
				log.Println(err.Error())
			}
		}
	}
}

func(p processLifeCycleRepository)  updateStatus(data v1.ProcessLifeCycleEvent,status string) error {
	filter := bson.M{
		"$and": []bson.M{
			{"process_id": data.ProcessId},
			{"step": data.Step},
		},
	}
	data.Status=enums.PROCESS_STATUS(status)
	update := bson.M{
		"$set": data,
	}
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	coll := p.manager.Db.Collection(ProcessLifeCycleCollection)
	err := coll.FindOneAndUpdate(p.manager.Ctx, filter, update, &opt)
	if err.Err() != nil {
		log.Println("[ERROR]", err.Err())
		return err.Err()
	}

	return nil
}
func(p processLifeCycleRepository)  update(data v1.ProcessLifeCycleEvent) error {
	filter := bson.M{
		"$and": []bson.M{
			{"process_id": data.ProcessId},
			{"step": data.Step},
		},
	}
	update := bson.M{
		"$set": data,
	}
	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}
	coll := p.manager.Db.Collection(ProcessLifeCycleCollection)
	err := coll.FindOneAndUpdate(p.manager.Ctx, filter, update, &opt)
	if err.Err() != nil {
		log.Println("[ERROR]", err.Err())
		return err.Err()
	}

	return nil
}
func (p processLifeCycleRepository) GetByProcessIdAndStep(processId, step string) *v1.ProcessLifeCycleEvent{
	query:=bson.M{
		"$and": []bson.M{
			{"process_id": processId},
			{"step": step},
		},
	}

	temp := new(v1.ProcessLifeCycleEvent)
	coll := p.manager.Db.Collection(ProcessLifeCycleCollection)
	result := coll.FindOne(p.manager.Ctx, query)
	err := result.Decode(temp)
	if err != nil {
		log.Println("[ERROR]", err)
	}
	if temp.ProcessId==""{
		return nil
	}
	return temp

}

func (p processLifeCycleRepository) GetByProcessId(processId string) []v1.ProcessLifeCycleEvent {
	query:=bson.M{
		"$and": []bson.M{
			{"process_id": processId},
		},
	}
	coll := p.manager.Db.Collection(ProcessLifeCycleCollection)

	curser, err := coll.Find(p.manager.Ctx, query)
	if err!=nil{
		log.Println(err.Error())
	}
	var results []v1.ProcessLifeCycleEvent
	for curser.Next(context.TODO()) {
		elemValue := new(v1.ProcessLifeCycleEvent)
		err := curser.Decode(elemValue)
		if err != nil {
			log.Println("[ERROR]", err)
			break
		}
		results= append(results, *elemValue)
	}
	return results
}

func NewProcessLifeCycleRepository(timeout int) repository.ProcessLifeCycleEventRepository {
	return &processLifeCycleRepository{
		manager: GetDmManager(),
		timeout: time.Duration(timeout),
	}

}
