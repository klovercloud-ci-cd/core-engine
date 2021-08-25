package main

import (
	"github.com/klovercloud-ci/config"
	"github.com/klovercloud-ci/repository/v1/mongo"
)

func main(){
	config.New()
	if config.Database == string(config.MONGO){
		mongo.GetDmManager()
	}
}
