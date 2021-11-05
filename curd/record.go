package curd

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"mongoStream/cmd"
	"sync"
	"time"
)

func Record(event bson.M, dest, dataMongo *MongoClient, db string) {
	var beforeData, documentId interface{}
	var collName string
	var err error
	operationType := event["operationType"]
	coll := event["ns"].(bson.M)["coll"]
	if coll != nil {
		collName = coll.(string)
	}

	documentKey, ok := event["documentKey"]
	if ok {
		documentId = documentKey.(bson.M)["_id"]
		beforeData, err = dest.GetDocument(db, collName, documentId)
		if err != nil && err != mongo.ErrNoDocuments {
			log.Printf("error get src data, db is %s, coll is %s, id %v, %s", db, collName, documentKey, err.Error())
		}
	} else {
		log.Printf("%s in %s has been dropped", collName, db)
		beforeData = fmt.Sprintf("%s in %s has been dropped", collName, db)
	}

	clusterTime := event["clusterTime"]

	record := map[string]interface{}{
		"op":     operationType,
		"ts":     clusterTime,
		"before": beforeData,
		"db": db,
		"coll": collName,
	}
	switch operationType {
	case "insert":
		full := event["fullDocument"]
		record["after"] = full
		if e := dest.InsertData(db, collName, format(full.(bson.M))); e != nil {
			log.Printf("error to insert new coll, db is %s, coll is %s, id %v, %s", db, collName, documentId, e.Error())
		}
	case "delete":
		record["after"] = nil
		if e := dest.DeleteCollData(db, collName, documentKey); e != nil {
			log.Printf("error to delete coll, db is %s, coll is %s, id %v, %s", db, collName, documentKey, e.Error())
		}
	case "drop":
		record["after"] = nil
		if e := dest.DropColl(db, collName); e != nil {
			log.Printf("error to drop coll: %s", e.Error())
		}
	case "update":
		full := event["fullDocument"]
		record["after"] = full
		for key, value := range event["updateDescription"].(bson.M) {
			record[key] = value
		}
		if e := dest.UpdateData(db, collName, documentKey, format(full.(bson.M))); e != nil {
			log.Printf("error to update coll, db is %s, coll is %s, id %v, %s", db, collName, documentKey, e.Error())
		}
	case "replace":
		full := event["fullDocument"]
		record["after"] = full
		if e := dest.ReplaceData(db, collName, documentKey, format(full.(bson.M))); e != nil {
			log.Printf("error to replace coll, db is %s, coll is %s, id %v, %s", db, collName, documentKey, e.Error())
		}
	case "dropDatabase":
		record["after"] = nil
		log.Printf("someone dropped the db %s", db)
	default:
		log.Printf("unknow operationType: %s", operationType)
	}

	// 将修改前的数据存储到一个新的库里
	go func() {
		if err := dataMongo.InsertData(cmd.DestDB, cmd.DestColl, record); err != nil {
			log.Printf("error to record history data, db is %s, coll is %s, id is %v, %s", db, collName, documentKey, err.Error())
		}
	}()
}

func Stream(src, dest, dataMongo *MongoClient, db string, wg *sync.WaitGroup) {
	defer wg.Done()
	srcDB := src.C.Database(db)
	streamOptions := options.ChangeStream().SetMaxAwaitTime(2 * time.Second)
	streamOptions.SetFullDocument(options.UpdateLookup)
	stream, err := srcDB.Watch(context.TODO(), mongo.Pipeline{}, streamOptions)
	if err != nil {
		log.Fatalf("error to start watch: %s", err.Error())
	}
	defer func() {
		if e := stream.Close(context.TODO()); e != nil {
			log.Fatalf("error to stop watch: %s", err.Error())
		}
	}()
	log.Printf("start to watch db %s...\n", db)
	for stream.Next(context.TODO()) {
		// A new event variable should be declared for each event.
		var event bson.M
		if err := stream.Decode(&event); err != nil {
			log.Fatalf("error to decode watch data: %s", err.Error())
		}
		go Record(event, dest, dataMongo, db)
	}
	if err := stream.Err(); err != nil {
		log.Fatalf("stream err: %s", err.Error())
	}
}

