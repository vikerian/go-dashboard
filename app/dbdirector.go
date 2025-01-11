package main

import "log/slog"

// types definitions
type Databases struct {
	Siri  SiriDB
	Redis RedisDB
	Mongo MongoDB
	log   *slog.Logger
}

func NewDBConnections(lg *slog.Logger, cfg *DBAppConfig) (dbs *Databases) {
	return &Databases{
		log: lg,
	}

}
