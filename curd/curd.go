package curd

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoClient struct {
	C *mongo.Client
}

func NewClient(mongoUri string) (*MongoClient, error) {
	clientOpts := options.Client().ApplyURI(mongoUri)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return nil, err
	}
	res := &MongoClient{
		C: client,
	}
	return res, nil
}

func (m *MongoClient) InsertData(db, coll string, data interface{}) (err error) {
	collection := m.C.Database(db).Collection(coll)
	_, err = collection.InsertOne(context.TODO(), data)
	return
}

func (m *MongoClient) DeleteCollData(db, coll string, id interface{}) (err error) {
	collection := m.C.Database(db).Collection(coll)
	_, err = collection.DeleteOne(context.TODO(), id)
	return
}

func (m *MongoClient) DropColl(db, coll string) (err error) {
	collection := m.C.Database(db).Collection(coll)
	err = collection.Drop(context.TODO())
	return
}

func (m *MongoClient) DropDatabase(db string) (err error) {
	d := m.C.Database(db)
	err = d.Drop(context.TODO())
	return
}

func (m *MongoClient) UpdateData(db, coll string, id, data interface{}) (err error) {
	opts := options.Update().SetUpsert(true)
	collection := m.C.Database(db).Collection(coll)
	update := bson.M{"$set": data}
	_, err = collection.UpdateOne(context.TODO(), id, update, opts)
	return
}

func (m *MongoClient) ReplaceData(db, coll string, id, data interface{}) (err error) {
	opts := options.Replace().SetUpsert(true)
	collection := m.C.Database(db).Collection(coll)
	_, err = collection.ReplaceOne(context.TODO(), id, data, opts)
	return
}

func (m *MongoClient) Disconnect() error {
	e := m.C.Disconnect(context.TODO())
	return e
}

func (m *MongoClient) GetDBList() (names []string, e error) {
	names, e = m.C.ListDatabaseNames(context.TODO(), bson.D{})
	return
}

func (m *MongoClient) GetDocument(db, coll string, id interface{}) (data bson.M, e error) {
	collection := m.C.Database(db).Collection(coll)
	e = collection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&data)
	return
}
