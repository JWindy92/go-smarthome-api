package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Device struct {
	Id    string `json:"Id"`
	Name  string `json:"Name"`
	State string `json:"State"`
}

var DummyDB []Device

func getAllDevices(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[GET] /devices")
	json.NewEncoder(w).Encode(DummyDB)
}

func getDeviceById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]

	fmt.Printf("[GET] /devices/%s", key)

	for _, device := range DummyDB {
		if device.Id == key {
			json.NewEncoder(w).Encode(device)
		}
	}
}

func createNewDevice(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[POST] /devices")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var device Device
	json.Unmarshal(reqBody, &device) // What is Unmarshal? What is '&' doing?

	DummyDB = append(DummyDB, device)

	json.NewEncoder(w).Encode(device)
}

func deleteArtice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Printf("[DELETE] /devices/%s", id)
	for idx, device := range DummyDB {
		if device.Id == id {
			DummyDB = append(DummyDB[:idx], DummyDB[idx+1:]...)
		}
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true) // What is StrictSlash?
	myRouter.HandleFunc("/devices", createNewDevice).Methods("POST")
	myRouter.HandleFunc("/devices", getAllDevices)
	myRouter.HandleFunc("/devices/{id}", deleteArtice).Methods("DELETE")
	myRouter.HandleFunc("/devices/{id}", getDeviceById).Methods("GET")

	log.Fatal(http.ListenAndServe(":5000", myRouter))
}

func main() {
	DummyDB = []Device{
		{Id: "1", Name: "device-1", State: "0"},
		{Id: "2", Name: "device-2", State: "0"},
	}
	handleRequests()
}
