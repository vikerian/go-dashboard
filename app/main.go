package main

import (
	"fmt"
	"io"
	"log"
	"os"
	//	"github.com/k0kubun/pp"
	//	"github.com/k0kubun/pp"
)

// types definitions
type Databases struct {
	Siri  SiriDB
	Redis RedisDB
	Mongo MongoDB
}

// global variables
var applog *log.Logger
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
	applog = log.New(mw, "", log.LstdFlags|log.Lmicroseconds)
}

// main function
func main() {
	// first some welcome (and defered goodbye :P)
	applog.Printf("Application starting up from binary %s...", os.Args[0])

	// read configuration
	config, err := ReadConfig("config.json")
	if err != nil {
		applog.Fatalf("Fatal error: %v", err)
	}

	// construct DSN for databases and connect to databases -> save connections to global dbs var
	redisDSN := fmt.Sprintf("redis://%s:%s@%s:%d/%d", config.Redis.Username, config.Redis.Password, config.Redis.Host, config.Redis.Port, config.Redis.DBIndex)
	siriDSN := fmt.Sprintf("siridb://%s:%s@%s:%d/%s", config.SiriDB.Username, config.SiriDB.Password, config.SiriDB.Host, config.SiriDB.Port, config.SiriDB.DBName)
	mongoDSN := MongoDBCreateDSN(config.Mongo.Username, config.Mongo.Password, config.Mongo.Host, config.Mongo.Port, config.Mongo.DBName)

	applog.Printf("Connecting to mongodb backend database...")
	dbs.Mongo, err = NewMongoConnection(mongoDSN)
	if err != nil {
		errstr := fmt.Sprintf("Error on connection to mongodb: %v", err)
		applog.Fatal(errstr)
	}
	applog.Printf("Connection complete...")

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

	// print instances of connections :)
	applog.Printf("Mongo DB conection instance: %v\n", dbs.Mongo)
	//pp.Print(dbs.Mongo)

	applog.Printf("Redis DB connection instance: %v\n", dbs.Redis)
	//pp.Print(dbs.Redis)

	applog.Printf("Siri DB connection instance: %v\n", dbs.Siri)
	//pp.Print(dbs.Siri)

	// end of main
	applog.Printf("Application binary %s ending...", os.Args[0])

	// close connections to db
	defer dbs.Redis.Close()
	defer dbs.Mongo.Close()
	defer dbs.Siri.Close()
}
