package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	zap "github.com/JWindy92/go-smarthome-api/pkg/logwrapper"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Zap = zap.NewLogger("main.go")

// var MqttClient = mqtt_utils.MqttInit()

// TODO: the main functions called by the handlers should be run concurrently using goroutines
func allDeviceHandler(w http.ResponseWriter, r *http.Request) {
	Zap.Logger.Infow(
		"Handling Req",
		"method", "GET",
		"route", "/devices",
	)
	var result = dbutils.GetAllDevices()

	json.NewEncoder(w).Encode(result)
}

func getDeviceByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	Zap.Logger.Infow(
		"Handling Req",
		"method", "GET",
		"route", "/devices/{id}",
		"id", id,
	)

	var result = dbutils.GetDeviceById(id)
	json.NewEncoder(w).Encode(result)
}

func getDevicesByTypeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	device_type := vars["type"]
	Zap.Logger.Infow(
		"Handling Req",
		"method", "GET",
		"route", "devices/type/{type}",
		"type", device_type,
	)

	var result = dbutils.GetDevicesOfType(device_type)
	json.NewEncoder(w).Encode(result)
}

func newDeviceHandler(w http.ResponseWriter, r *http.Request) {
	Zap.Logger.Infow(
		"Handling Req",
		"method", "POST",
		"route", "/devices",
	)

	var prim primitive.M
	// reqBody, _ := ioutil.ReadAll(r.Body)
	_ = json.NewDecoder(r.Body).Decode(&prim)

	var result = dbutils.CreateNewDevice(prim)
	json.NewEncoder(w).Encode(result)
}

func deleteDeviceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	Zap.Logger.Infow(
		"Handling Req",
		"method", "DELETE",
		"route", "/devices/{id}", //? Can I format this value to contain the actual Id?
		"id", id,
	)

	dbutils.DeleteDevice(id)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true) // What is StrictSlash?
	myRouter.HandleFunc("/devices", newDeviceHandler).Methods("POST")
	myRouter.HandleFunc("/devices", allDeviceHandler)
	myRouter.HandleFunc("/devices/{id}", deleteDeviceHandler).Methods("DELETE")
	myRouter.HandleFunc("/devices/{id}", getDeviceByIdHandler).Methods("GET")

	myRouter.HandleFunc("/devices/type/{type}", getDevicesByTypeHandler).Methods("GET")
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
