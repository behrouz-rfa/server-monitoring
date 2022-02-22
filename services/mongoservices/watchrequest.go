package mongoservices

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"server-monitoring/domain/requests"
	"server-monitoring/shared/database"
)

// this is test for watch request on
// mongo db
func Run() {
	coll := database.Mongo.Database("monitoring").Collection("requests")
	pipLine := mongo.Pipeline{bson.D{{"$match", bson.D{{"operationType", "insert"}}}}}
	cs, err := coll.Watch(context.TODO(), pipLine)
	if err != nil {
		fmt.Println(err)
	}
	defer cs.Close(context.TODO())
	fmt.Println("Waiting for chane envents insert somthing in collection")

	for cs.Next(context.TODO()) {
		var event bson.M
		if err := cs.Decode(&event); err != nil {
			panic(err)
		}
		output, err := json.MarshalIndent(event["fullDocument"], "", "    ")

		if err != nil {
			fmt.Println(err)
		}
		var request requests.Request
		err = json.Unmarshal(output, &request)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("%s\n", output)
	}

	if err := cs.Err(); err != nil {
		fmt.Println(err)
	}

}
