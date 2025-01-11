package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	//	"github.com/k0kubun/pp"
	//	"github.com/k0kubun/pp"
)

// global variables
var applog *slog.Logger
var dbs Databases

// we use init function to setup logging :P
// default log file: application.log
func init() {
	lf, err := os.OpenFile("application.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0640)
	if err != nil {
		panic(err)
	}
	// multiwriter
	mw := io.MultiWriter(os.Stdout, lf)
	// our logger setup
	applog = slog.New(slog.NewJSONHandler(mw, nil))
}

// main function
func main() {
	// first some welcome (and defered goodbye :P)
	applog.Info("Application starting up from binary %s...", os.Args[0])

	// read configuration
	config, err := ReadConfig("config.json")
	if err != nil {
		applog.Error("Fatal error: %v", err)
		panic(err)
	}

	// construct DSN for databases and connect to databases -> save connections to global dbs var
	//redisDSN := fmt.Sprintf("redis://%s:%s@%s:%d/%d", config.Redis.Username, config.Redis.Password, config.Redis.Host, config.Redis.Port, config.Redis.DBIndex)
	//siriDSN := fmt.Sprintf("siridb://%s:%s@%s:%d/%s", config.SiriDB.Username, config.SiriDB.Password, config.SiriDB.Host, config.SiriDB.Port, config.SiriDB.DBName)
	mongoDSN := MongoDBCreateDSN(config.Mongo.Username, config.Mongo.Password, config.Mongo.Host, config.Mongo.Port, config.Mongo.DBName)

	applog.Info("Connecting to mongodb backend database...")
	dbs.Mongo, err = NewMongoConnection(mongoDSN)
	if err != nil {
		errstr := fmt.Sprintf("Error on connection to mongodb: %v", err)
		applog.Error(errstr)
		panic(err)
	}
	applog.Info("Connection complete...")
	/*
		applog.Printf("Connecting to redis database engine...")
		dbs.Redis, err = NewRedisConnection(redisDSN)
		if err != nil {
			errstr := fmt.Sprintf("Error on connectio to redis: %v", err)
			applog.Fatal(errstr)
		}
		applog.Printf("Connection complete...")

		//siridb := SiriDB()
		applog.Printf("Connecting to SiriDB timescale database...")
		dbs.Siri, err = NewSiriDBConnection(siriDSN)
		if err != nil {
			errstr := fmt.Sprintf("Error on connection to siridb : %v", err)
			applog.Fatal(errstr)
		}
		applog.Printf("Connection complete")
	*/
	// print instances of connections :)
	applog.Debug("Mongo DB conection instance: %v\n", dbs.Mongo)
	//pp.Print(dbs.Mongo)

	//applog.Printf("Redis DB connection instance: %v\n", dbs.Redis)
	//pp.Print(dbs.Redis)

	//applog.Printf("Siri DB connection instance: %v\n", dbs.Siri)
	//pp.Print(dbs.Siri)

	// end of main
	applog.Info("Application binary %s ending...", os.Args[0])

	// close connections to db
	//defer dbs.Redis.Close()
	defer dbs.Mongo.Close()
	//defer dbs.Siri.Close()
}
