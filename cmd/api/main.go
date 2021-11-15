package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	zap "github.com/JWindy92/go-smarthome-api/pkg/logwrapper"
	"github.com/gorilla/mux"
)

var Zap = zap.NewLogger()

func allDeviceHandler(w http.ResponseWriter, r *http.Request) {
	Zap.Logger.Infof(
		"Handling Req",
		"Method", "GET",
		"Route", "/devices",
	)
	var result = dbutils.GetAllDevices()

	json.NewEncoder(w).Encode(result)
}

func getDeviceByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	Zap.Logger.Infof(
		"Handling Req",
		"Method", "GET",
		"Route", "/devices/{id}",
		"Id", id,
	)

	var result = dbutils.GetDeviceById(id)
	json.NewEncoder(w).Encode(result)
}

func newDeviceHandler(w http.ResponseWriter, r *http.Request) {
	Zap.Logger.Infof(
		"Handling Req",
		"Method", "POST",
		"Route", "/devices",
	)
	reqBody, _ := ioutil.ReadAll(r.Body)

	var result = dbutils.CreateNewDevice(reqBody)
	json.NewEncoder(w).Encode(result)
}

func deleteDeviceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	Zap.Logger.Infof(
		"Handling Req",
		"Method", "DELETE",
		"Route", "/devices/{id}", //? Can I format this value to contain the actual Id?
		"Id", id,
	)

	dbutils.DeleteDevice(id)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true) // What is StrictSlash?
	myRouter.HandleFunc("/devices", newDeviceHandler).Methods("POST")
	myRouter.HandleFunc("/devices", allDeviceHandler)
	myRouter.HandleFunc("/devices/{id}", deleteDeviceHandler).Methods("DELETE")
	myRouter.HandleFunc("/devices/{id}", getDeviceByIdHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000", myRouter))
}

func main() {

	defer Zap.Logger.Sync() // TODO: Find out if this single call in main is sufficent

	// TODO: these values should be populated by variables
	Zap.Logger.Infow(
		"Starting API server",
		"host", "localhost",
		"port", 5000,
	)

	Zap.Logger.Infof("Ready to accept requests")
	handleRequests()
}
