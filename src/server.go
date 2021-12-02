package server

import (
	"duke/init/src/cache"
	"duke/init/src/database"
	util "duke/init/src/helpers"
	"duke/init/src/router"
	"fmt"
	login "github.com/Albinzr/duke_login_module"
	"github.com/Albinzr/duke_product_module/product"
	"log"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"runtime"
)

//Message :- simple type for message callback
type Message func(message string)

var env = util.EnvConfig()
var path, _ = filepath.Abs("./store")

var cacheConfig = &cache.Config{
	Host:       "localhost",
	Port:       "6379",
	DB:         0,
	Password:   "",
	MaxRetries: 3,
}

var dbConfig = &database.Config{
	URL:          env.MongoURL,
	DatabaseName: env.DatabaseName,
}

//Start :- server start function
func Start() {
	go cacheConfig.Init()
	dbConfig.Init()

	loginConfig := &login.Config{
		Database:               dbConfig.Database,
		CollectionName:         "user",
		Aud:                    env.Aud,
		Iss:                    env.Iss,
		ForgotPasswordCallback: resetPasswordEmail,
		SecretKey:              env.SecretKey,
	}
	fmt.Println(env.SecretKey, "------->")
	go loginConfig.Init()
	productConfig := product.Config{
		Database:       dbConfig.Database,
		CollectionName: "product",
	}
	go productConfig.Init()
	go router.Init()
	runServer()
}

func resetPasswordEmail(emailId string, url string) {
	util.LogInfo(emailId, "----->", url)
}

func runServer() {
	util.LogInfo("Listing for requests at http://localhost:1000/")
	runtime.GOMAXPROCS(runtime.NumCPU())
	//go util.PrintMemUsageWithTimer()
	log.Fatal(http.ListenAndServe(":1000", nil))
}
