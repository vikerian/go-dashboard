package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	//"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
	//"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// "go.mongodb.org/mongo-driver/mongo"

type MongoCon struct {
	URL        string
	CTX        context.Context
	Cancel     context.CancelFunc
	Options    *options.ClientOptions
	CLH        *mongo.Client
	Collection *mongo.Collection // aka "table"
}

type MongoDB interface {
	//NewMongoConnection(string) (*MongoCon, error)
	Create(string, interface{}) error
	Read(string) (interface{}, error)
	Delete(string) error
	Close() error
}

// ConstructDSN - create DSN for connection, mongo is quite picky on this
func MongoDBCreateDSN(username, password, host, port, database string) string {
	var mongoDSN string
	if username == "" || password == "" {
		mongoDSN = fmt.Sprintf("mongodb://%s:%s/%s", host, port, database)
	} else {
		//if username != "admin" {
		mongoDSN = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", username, password, host, port, database)
		//} else {
		//	mongoDSN = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", username, password, host, port, database)
	}
	return mongoDSN
}

// NewMongoConnection - constructor for our communications with Mongo
func NewMongoConnection(dsn string) (*MongoCon, error) {
	// create instance of MongoConnection
	ma := new(MongoCon)
	ma.URL = dsn
	ma.CTX, ma.Cancel = context.WithTimeout(context.Background(), 5*time.Second)
	ma.Options = options.Client().ApplyURI(dsn)
	clh, err := mongo.Connect(ma.CTX, ma.Options)
	ma.CLH = clh
	if err != nil {
		errstr := fmt.Sprintf("Error on setup connection to Mongo: %v", err)
		return nil, errors.New(errstr)
	}
	err = ma.CLH.Ping(ma.CTX, nil)
	if err != nil {
		errstr := fmt.Sprintf("Error on checkup connection to Mongo: %v", err)
		return nil, errors.New(errstr)
	}
	return ma, nil
}

// CreateVAL -> create value with specified key on ourcollection
func (mc *MongoCon) Create(key string, value interface{}) error {

	return nil
}

// ReadVAL -> read value specified by key
func (mc *MongoCon) Read(key string) (interface{}, error) {

	return nil, nil
}

func (mc *MongoCon) Delete(key string) error {

	return nil
}

func (mc *MongoCon) Close() error {
	return mc.CLH.Disconnect(mc.CTX)
}
