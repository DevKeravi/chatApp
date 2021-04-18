package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type dbHandler struct {
	db   *mongo.Client
	coll *mongo.Collection
}

var mgo dbHandler

func initDB() {
	mgo.db, _ = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:10102"))
	err := mgo.db.Connect(context.TODO())
	if err != nil {
		log.Fatal("mongoDB Connect() error : ", err)
	}
	log.Println("Connected mongoDB....")
	mgo.coll = mgo.db.Database("chatAppLog").Collection("message")

}
func CloseDB() {
	err := mgo.db.Disconnect(context.TODO())
	if err != nil {
		log.Println("mongoDB Disconnect() error : ", err)
	}
}

func writeDB(msg message) {
	_, err := mgo.coll.InsertOne(context.TODO(), msg)
	if err != nil {
		log.Println(err)
		return
	}
}

func getDB() []message {

	list := []message{}
	filter := bson.D{{}}

	cursor, err := mgo.coll.Find(context.TODO(), filter.Map())
	if err != nil {
		log.Println("mongoDB Find() error: ", err)
	}

	for cursor.Next(context.TODO()) {
		var temp message
		err := cursor.Decode(&temp)
		if err != nil {
			log.Println("mongoDB cursor Decode error: ", err)
		}
		list = append(list, temp)
	}

	return list
}
