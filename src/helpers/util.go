package util

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

//Config :- env struct
type Config struct {
	Port string
	MongoURL     string
	DatabaseName string
}

//LogError :- common function for loging error
func LogError(message string, errorData error) {
	if errorData != nil {
		log.Errorln("Error : ", message)
		return
	}
}

//LogInfo :- common func for loging info
func LogInfo(args ...interface{}) {
	log.Info(args)
}

//LogFatal :- common func for fatal error
func LogFatal(args ...interface{}) {
	log.Fatal(args)
}

//LogDebug :- common debug logger
func LogDebug(args ...interface{}) {
	log.Debug(args)
}

//LoadEnvConfig :- for loading config files
func LoadEnvConfig() *Config {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var err error

	key := flag.String("env", "development", "")
	flag.Parse()
	LogInfo("env:", *key)
	if *key == "production" {
		log.SetFormatter(&log.TextFormatter{})
		err = godotenv.Load("./production.env")
	} else {
		err = godotenv.Load("./local.env")
		log.SetFormatter(&log.TextFormatter{})
	}

	if err != nil {
		LogFatal("cannot load config file", err)
	}

	config := new(Config)
	config.Port = os.Getenv("PORT")
	config.MongoURL = os.Getenv("MONGO_URL")
	config.DatabaseName = os.Getenv("DATABASE_NAME")
	return config
}

//PrintMemUsage -test
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
	fmt.Printf("\tMemory Freed = %v\n", bToMb(m.Frees))

	runtime.GC()
	debug.FreeOSMemory()
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
func PrintMemUsageWithTimer(){
	for now := range time.Tick(time.Minute) {
		fmt.Println(now)
		PrintMemUsage()
	}
}