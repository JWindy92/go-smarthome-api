package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func allDeviceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[GET] /devices")
	var result = getAllDevices()

	json.NewEncoder(w).Encode(result)
}

func getDeviceByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Printf("[GET] /devices/%s", id)

	var result = getDeviceById(id)
	json.NewEncoder(w).Encode(result)
}

func newDeviceHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[POST] /devices")
	reqBody, _ := ioutil.ReadAll(r.Body)

	var result = createNewDevice(reqBody)
	json.NewEncoder(w).Encode(result)
}

func deleteDeviceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Printf("[DELETE] /devices/%s", id)
	deleteDevice(id)
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
	DummyDB = []Device{
		{Id: "1", Name: "device-1", State: "0"},
		{Id: "2", Name: "device-2", State: "0"},
	}
	handleRequests()
}
