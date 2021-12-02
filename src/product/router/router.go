package router

import (
	util "duke/init/src/helpers"
	"duke/init/src/product/Config"
	"duke/init/src/product/database"
	"fmt"
	"net/http"
	"strconv"
)

type Config ProductConfig.Config

var dbConfig *database.Config

func (c *Config) Init() {

	dbConfig = (*database.Config)(c)
	dbConfig.Init()

	http.HandleFunc("/create", c.createHandler)
	http.HandleFunc("/update", c.updateHandler)
	http.HandleFunc("/delete", c.deleteHandler)
	//
	http.HandleFunc("/listAllProduct", c.listAllProductHandler)
	http.HandleFunc("/product", c.listProductHandler)
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
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			resp := util.ErrorResponse("unable to save data", "data not saved", err)
			_, _ = w.Write(resp)
			return
		}
	fmt.Print(objId, err)

	w.WriteHeader(http.StatusOK)
	resp := util.SuccessResponse(`{"productId":"` + objId.Hex() +`"}`)
	_, _ = w.Write(resp)
	return
}

func (c *Config) updateHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := req.ParseForm()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := util.ErrorResponse("invalid param", "param not valid", err)
		_, _ = w.Write(resp)
		return
	}

	isUpadted, err := dbConfig.Update(req.Form)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := util.ErrorResponse("unable to save data", "data not updated", err)
		_, _ = w.Write(resp)
		return
	}
	fmt.Print(isUpadted, err)
	if isUpadted{
		w.WriteHeader(http.StatusOK)
		resp := util.SuccessResponse("null")
		_, _ = w.Write(resp)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	resp := util.ErrorResponse("unable to update ", "data not updated", err)
	_, _ = w.Write(resp)
	return
}

func (c *Config) deleteHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := req.ParseForm()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := util.ErrorResponse("invalid param", "param not valid", err)
		_, _ = w.Write(resp)
		return
	}

	for key, val := range  req.Form{
		print(key,val)
	}

	id := req.Form.Get("_id")
	fmt.Println(id,"----><-----")
	isRemoved, err := dbConfig.Delete(id)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := util.ErrorResponse("unable to delete product", "product not deleted from db", err)
		_, _ = w.Write(resp)
		return
	}
	fmt.Print(isRemoved, err)

	w.WriteHeader(http.StatusOK)
	resp := util.SuccessResponse(`{"removed":"` + strconv.FormatBool(isRemoved) +`"}`)
	_, _ = w.Write(resp)
	return
}

func (c *Config) listAllProductHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := req.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := util.ErrorResponse("invalid param", "param not valid", err)
		_, _ = w.Write(resp)
		return
	}

	limit, err := strconv.ParseInt(req.Form.Get("limit"), 10, 64)

	if err != nil{
		limit = 10
	}

	offset, err := strconv.ParseInt(req.Form.Get("offset"), 10, 64)
	if err != nil{
		offset = 0
	}

	products, err := dbConfig.FindAllProduct(&offset,&limit)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := util.ErrorResponse("invalid param", "param not valid", err)
		_, _ = w.Write(resp)
		return
	}

	fmt.Println(products)
	w.WriteHeader(http.StatusOK)
	resp := util.SuccessResponseWithInterface(products)
	_, _ = w.Write(resp)
	return

}

func (c *Config) listProductHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := req.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := util.ErrorResponse("invalid param", "param not valid", err)
		_, _ = w.Write(resp)
		return
	}
	id := req.Form.Get("_id")

	product ,err := dbConfig.FindSingleProduct(id)


	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := util.ErrorResponse("invalid param", "param not valid", err)
		_, _ = w.Write(resp)
		return
	}

	fmt.Println(product)
	w.WriteHeader(http.StatusOK)
	resp := util.SuccessResponseWithInterface(product)
	_, _ = w.Write(resp)
	return
}
