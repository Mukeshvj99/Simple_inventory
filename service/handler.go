package service

import (
	// "github.com/mukesh/Simple_Inventory/main"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/mukesh/Simple_Inventory/persistence"
)

func CloseDb() {
	persistence.CloseDb()
}
func GetConnection(url string) {
	err := persistence.MakeConnection(url)
	if err != nil {
		log.Fatal("cannot make connection to database")
	}

}

func SetTablename(table string) {
	persistence.TableData(table)
	persistence.InitializeQueries()
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {

	var p []persistence.Product
	body := json.NewDecoder(r.Body)

	decodeerr := body.Decode(&p)
	if decodeerr != nil {
		Response(w, 400, formatResponse("Invalid Request Body"))
		return
	}

	create, err := persistence.CreateProductData(p)

	if create == 0 {
		Response(w, 400, formatResponse("The product Id Already exist in the DB"))
		return
	}

	Response(w, 200, map[string]interface{}{"msg": "successfully Created", "failedID": err, "success": create})
	return

}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		Response(w, 400, formatResponse("Invalid product id"))
		return
	}

	var p persistence.Product
	body := json.NewDecoder(r.Body)

	decodeerr := body.Decode(&p)
	if decodeerr != nil {
		Response(w, 400, formatResponse("Invalid Request Body"))
		return
	}

	update, err := persistence.UpdataProductData(id, p)

	if update == 0 {
		Response(w, 400, formatResponse("The product Id doesn't exist in the DB"))
		return
	}
	if err != nil && err.Error() == "Error in deleting rows" {
		Response(w, 500, formatResponse("Internal db issue"))
		return
	}

	Response(w, 200, map[string]string{"msg": "successfully Updated"})
	return
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		Response(w, 400, formatResponse("Invalid product id"))
		return
	}

	delete, err := persistence.DeleteProductData(id)

	if delete == 0 {
		Response(w, 400, formatResponse("The product Id doesn't exist in the DB"))
		return
	}
	if err != nil && err.Error() == "Error in deleting rows" {
		Response(w, 500, formatResponse("Internal db issue"))
		return
	}

	Response(w, 200, map[string]string{"msg": "successfully deleted"})
	return
}

func GetProduct(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	productid, err := strconv.Atoi(vars["id"])
	if err != nil {
		Response(w, 400, formatResponse("Invalid product id"))
		return
	}

	var p persistence.Product

	rowerr := p.GetProductData(int64(productid))

	if rowerr != nil && rowerr.Error() == "Error in fetching rows" {
		Response(w, 500, formatResponse("Internal db issue"))
		return
	}
	if rowerr != nil {
		log.Println(rowerr.Error())
		Response(w, 400, formatResponse("Product Doesn't Exist in DB"))
		return
	}

	Response(w, 200, p)
	return

}

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	start, err := strconv.Atoi(r.FormValue("start"))
	count, err1 := strconv.Atoi(r.FormValue("count"))

	log.Println("start", start, "count", count)
	if err != nil || err1 != nil {
		Response(w, 400, formatResponse("Invalid start or end count"))
		return
	}

	if start < 0 || count < 0 {
		Response(w, 400, formatResponse("Invalid start or end count"))
		return
	}

	products, err := persistence.GetAllProductsData(count, start)

	if err != nil && err.Error() == "Error in fetching rows" {
		log.Println("error", err)
		Response(w, 500, formatResponse("Internal db issue"))
		return
	}

	Response(w, 200, products)
	return

}

func formatResponse(message string) map[string]interface{} {
	res := make(map[string]interface{})
	res["error"] = message
	return res
}
func Response(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
