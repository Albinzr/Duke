package router

import (
	util "duke/init/src/helpers"
	"duke/init/src/product/Config"
	"duke/init/src/product/database"
	"fmt"
	"net/http"
)

type Config ProductConfig.Config

var dbConfig *database.Config

func (c *Config) Init() {

	dbConfig = (*database.Config)(c)
	dbConfig.Init()

	http.HandleFunc("/create", c.createHandler)
	http.HandleFunc("/update", c.handler)
	http.HandleFunc("/delete", c.handler)
	//
	http.HandleFunc("/listAllProduct", c.handler)
	http.HandleFunc("/listProduct", c.handler)
}

func (c *Config) createHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := req.ParseForm()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := util.ErrorResponse("invalid param", "param not valid", err)
		_, _ = w.Write(resp)
		return
	}

	objId, err := dbConfig.Create(req.Form)

	fmt.Print(objId, err)

	w.WriteHeader(http.StatusOK)
	resp := util.SuccessResponse(`{"productId":1}`)
	_, _ = w.Write(resp)
	return
}

func (c *Config) handler(w http.ResponseWriter, req *http.Request) {

}
