package cmd

import "flag"

var (
	SrcMongoUri  = ""
	DestMongoUri = ""
	DBs          = ""
	DataMongoUri = ""
	DestDB       = ""
	DestColl     = ""
)

func init() {
	flag.StringVar(&SrcMongoUri, "src-mongo-uri", "mongodb://localhost:27017", "which mongo to watch，format: mongodb://user:pass@addr:port or mongodb://addr:port/?replicaSet=rs0")
	flag.StringVar(&DestMongoUri, "dest-mongo-uri", "mongodb://localhost:27017", "which mongo to sync the data from src mongo，format: mongodb://user:pass@addr:port or mongodb://addr:port/?replicaSet=rs0")
	flag.StringVar(&DBs, "dbs", "all", "which dbs to watch, as a comma separated list, default is all dbs")
	flag.StringVar(&DataMongoUri, "data-mongo-uri", SrcMongoUri, "which mongo to store the change history, default is src-mongo-uri")
	flag.StringVar(&DestDB, "dest-db", "records", "which db to store the change steam data")
	flag.StringVar(&DestColl, "dest-coll", "mongo", "which coll to store the change steam data")
	flag.Parse()
}