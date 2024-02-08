package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	//"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

func main() {

	configpath := "../config.yaml"
	env := *flag.String("env", "local", "if you are running on multiple environments")
	port := *flag.String("port", "8080", "mention the port number to run the application")
	databaseconfig, err := Loadconfig(configpath, env)

	if err != nil {
		panic("cannot load the configuration check it")
	}
	log.Println("Our database config", databaseconfig)

	inventory := Inventory_App{}

	inventory.Initialize(databaseconfig.Username, databaseconfig.Password, databaseconfig.Database, databaseconfig.Tablename)

	address := "127.0.0.1:" + port
	inventory.Run(address)

	log.Println("Inventory Application started")

}

func Loadconfig(path, env string) (Database, error) {

	var databaseconfig Database
	var config map[string]Database
	data, err := ioutil.ReadFile(path)

	if err != nil {
		return databaseconfig, err
	}

	if err := yaml.Unmarshal(data, &config); err != nil {

		fmt.Println("errror", err)
		return Database{}, err
	}
	fmt.Println("config", config)

	return config[env], nil
}
