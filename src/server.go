package server

import (
	"duke/init/src/cache"
	"duke/init/src/database"
	util "duke/init/src/helpers"
	"duke/init/src/router"
	"log"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
)

//Message :- simple type for message callback
type Message func(message string)

var env = util.LoadEnvConfig()
var path, _ = filepath.Abs("./store")

var cacheConfig = &cache.Config{
	Host: "localhost",
	Port: "6379",
	DB: 0,
	Password: "",
	MaxRetries: 3,
}

var dbConfig = &database.Config{
	URL:          env.MongoURL,
	DatabaseName: env.DatabaseName,
}


//Start :- server start function
func Start() {
	go cacheConfig.Init()
	go router.Init()
	err := dbConfig.Init()

	if err != nil {
		util.LogError("Database connection issue", err)
		return
	}else{
		util.LogInfo("Database connected")
	}


	log.Println("Listing for requests at http://localhost:1000/")
	go util.PrintMemUsageWithTimer()
	log.Fatal(http.ListenAndServe(":1000", nil))
}

