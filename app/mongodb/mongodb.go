package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDatabase *mongo.Database

// OpenMongoDB to open MongoDB connection
func OpenMongoDB(dbhost string) (*mongo.Client, error) {
	var (
		mc  *mongo.Client
		err error
	)
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()
	mc, err = mongo.Connect(ctx, options.Client().ApplyURI(dbhost))
	if err == nil {
		err = mc.Ping(ctx, nil)
	}

	return mc, err
}

// SetInstance init mongo database
func SetInstance(d *mongo.Database) {
	mongoDatabase = d
}

// GetInstance ...
func GetInstance() *mongo.Database {
	return mongoDatabase
}

// ToDoc ...
func ToDoc(v interface{}) (interface{}, error) {
	var doc interface{}

	data, err := bson.Marshal(v)
	if err == nil {
		err = bson.Unmarshal(data, &doc)
	}
	return doc, err
}
