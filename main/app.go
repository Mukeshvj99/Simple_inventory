package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	//	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mukesh/Simple_Inventory/service"
)

type Inventory_App struct {
	Router *mux.Router
}

type Database struct {
	Username  string
	Password  string
	Database  string
	Tablename string
}

func (a *Inventory_App) Run(address string) {

	server := &http.Server{
		Handler: a.Router,
		Addr:    address,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())

	defer func() {
		server.Close()
		log.Println("Closing the  application ..... ")
	}()
}

func (a *Inventory_App) Initialize(username, pass, dbname, tablename string) {

	database_url := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?", username, pass, dbname)

	service.GetConnection(database_url)
	fmt.Println("tablename", tablename)
	service.SetTablename(tablename)
	a.Router = mux.NewRouter()
	a.Routes()
}

func (a *Inventory_App) Routes() {

	a.Router.HandleFunc("/product/{id:[0-9]+}", service.GetProduct).Methods("GET")
	a.Router.HandleFunc("/products", service.GetAllProducts).Methods("GET")
	a.Router.HandleFunc("/products", service.CreateProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", service.UpdateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", service.DeleteProduct).Methods("DELETE")
}
