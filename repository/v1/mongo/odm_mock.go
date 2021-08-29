package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)


func GetMockDmManager() *dmManager {
	onceDmManager.Do(func() {
		log.Println("[INFO] Starting Initializing Singleton DB Manager")
		singletonDmManager = &dmManager{}
		singletonDmManager.initMockMongoConnection()
	})
	return singletonDmManager
}

func (dm *dmManager) initMockMongoConnection() {
	ctx := context.Background()
	dm.Ctx = ctx
	clientOpts := options.Client().ApplyURI( "mongodb://" + os.Getenv("MONGO_USERNAME") + ":" + os.Getenv("MONGO_PASSWORD") + "@" + os.Getenv("MONGO_SERVER") + ":" +  os.Getenv("MONGO_PORT"))
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Println("[ERROR] DB Connection error:", err.Error())
		return
	}

	db := client.Database(os.Getenv("DATABASE_NAME"))
	dm.Db = db

	log.Println("[INFO] Initialized Singleton DB Manager")
}


