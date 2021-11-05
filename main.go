package main

import (
	"log"
	"mongoStream/cmd"
	"mongoStream/curd"
	"strings"
	"sync"
)

func main() {
	srcClient, e := curd.NewClient(cmd.SrcMongoUri)
	if e != nil {
		log.Fatal("error to create srcClient " + e.Error())
	}
	destClient, e := curd.NewClient(cmd.DestMongoUri)
	if e != nil {
		log.Fatal("error to create destClient " + e.Error())
	}
	dataMongoClient, e := curd.NewClient(cmd.DataMongoUri)
	if e != nil {
		log.Fatal("error to create dataMongoClient " + e.Error())
	}

	defer func() {
		if e := srcClient.Disconnect(); e != nil {
			log.Fatal("error to close srcClient " + e.Error())
		}
		if e := destClient.Disconnect(); e != nil {
			log.Fatal("error to close destClient " + e.Error())
		}
		if e := dataMongoClient.Disconnect(); e != nil {
			log.Fatal("error to close destClient " + e.Error())
		}
	}()

	var dbs []string

	if strings.ToLower(strings.TrimSpace(cmd.DBs)) == "all" {
		names, e := srcClient.GetDBList()
		if e != nil {
			log.Fatalf("can't list db names: %s", e.Error())
		}
		dbs = names
	} else {
		dbs = strings.Split(cmd.DBs, ",")
	}

	wg := new(sync.WaitGroup)
	for key := range dbs {
		wg.Add(1)
		go curd.Stream(srcClient, destClient, dataMongoClient, strings.TrimSpace(dbs[key]), wg)
	}
	wg.Wait()
}
