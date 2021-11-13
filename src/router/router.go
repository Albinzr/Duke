package router

import (
	middleware "duke/init/src/middlewares"
	"net/http"
)

func Init() {
	http.Handle("/profile", middleware.IsAuthorized(profileHandler))
}

var profileHandler = func(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"profileHandler"}`))
}
