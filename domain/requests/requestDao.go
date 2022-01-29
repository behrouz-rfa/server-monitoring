package requests

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"server-monitoring/shared/database"
	loogers "server-monitoring/utils/looger"
)

func (r *Request) Find() ([]Request, error) {
	var results []Request
	cur, err := database.Mongo.Database("monitoring").Collection("requests").Find(context.Background(), bson.D{{}})
	if err != nil {
		return results, err
	}

	for cur.Next(context.TODO()) {
		var elem Request
		err := cur.Decode(&elem)
		if err != nil {
			continue
		}

		results = append(results, elem)
	}

	return results, nil
}

func (r *Request)  InsertConsoleLog() error {

	coll := database.Mongo.Database("monitoring").Collection("requests")

	//var nodes []interface{}

	//test1 := map[string]interface{}{
	//	"src_ip":    item.SrcIp,
	//	"src_mac":   item.SrcMac,
	//	"dst_ip":    item.DstIp,
	//	"dst_mac":   item.DstMac,
	//	"protocol":  item.Protocol,
	//	"timestamp": item.Timestamp,
	//}
	nod := bson.D{
		{"src_addr", (r.SrcAddr)},
		{"src_port", r.SrcPort},
		{"dst_addr", r.DstAddr},
		{"dst_port", r.DstPort},
		{"method", r.Method},
		{"ts", r.Ts},
		{"status_code", r.StatusCode},
		{"content_length", r.ContentLength},
		{"url", r.Url},
		{"user_agent", r.UserAgent},
		{"body", r.Body},
	}
	_, err := coll.InsertOne(context.TODO(), nod)
	if err != nil {
		loogers.Error("insert node", err)
		//return err
	}
	//nod := bson.D{
	//	{"src_ip", item.SrcIp},
	//	{"src_mac", item.SrcMac},
	//	{"dst_ip", item.DstIp},
	//	{"dst_mac", item.DstMac},
	//	{"protocol", item.Protocol},
	//	{"timestamp", item.Timestamp}}
	//nodes = append(nodes, test1)

	//_, err := nodeCollection.InsertMany(context.TODO(), nodes)
	//if err != nil {
	//	loogers.Error("insert node", err)
	//	return err
	//}

	return nil
}
