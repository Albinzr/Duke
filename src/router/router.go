package router

import (
	"net/http"
)

func Init(){
	http.HandleFunc("/", helloHandler)
}

var helloHandler = func(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte(`{"status":"running"}`))
}