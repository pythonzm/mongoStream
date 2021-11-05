# mongoStream
跟踪mongo数据变化

```bash
$ go mod
$ go build
$ mongoStream -h
Usage of mongoStream:
  -data-mongo-uri string
        which mongo to store the change history, default is src-mongo-uri (default "mongodb://localhost:27017")
  -dbs string
        which dbs to watch, as a comma separated list, default is all dbs (default "all")
  -dest-coll string
        which coll to store the change steam data (default "mongo")
  -dest-db string
        which db to store the change steam data (default "records")
  -dest-mongo-uri string
        which mongo to sync the data from src mongo，format: mongodb://user:pass@addr:port or mongodb://addr:port/?replicaSet=rs0 (default "mongodb://localhost:27017")
  -src-mongo-uri string
        which mongo to watch，format: mongodb://user:pass@addr:port or mongodb://addr:port/?replicaSet=rs0 (default "mongodb://localhost:27017")
```

* -src-mongo-uri 指定数据源mongo
* -dest-mongo-uri 指定中间库
* -data-mongo-uri 指定存储变更数据的mongo
