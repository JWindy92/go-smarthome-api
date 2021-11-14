package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	"github.com/gorilla/mux"
)

func allDeviceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[GET] /devices")
	var result = dbutils.GetAllDevices()

	json.NewEncoder(w).Encode(result)
}

func getDeviceByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Printf("[GET] /devices/%s", id)

	var result = dbutils.GetDeviceById(id)
	json.NewEncoder(w).Encode(result)
}

func newDeviceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[POST] /devices")
	reqBody, _ := ioutil.ReadAll(r.Body)

	var result = dbutils.CreateNewDevice(reqBody)
	json.NewEncoder(w).Encode(result)
}

func deleteDeviceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Printf("[DELETE] /devices/%s", id)
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
	handleRequests()
}
