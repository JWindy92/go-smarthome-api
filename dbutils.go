package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Device struct {
	Id    string `json:"Id"`
	Name  string `json:"Name"`
	State string `json:"State"`
}

// TODO: dbutils should not have anything to do with handling the request, need to modify these functions
// to work independently of http requst functionality. Then the functions called by the request handlers
// should also be modified to consume the request/body and prepare it for the DB functions

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

func deleteDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	fmt.Printf("[DELETE] /devices/%s", id)
	for idx, device := range DummyDB {
		if device.Id == id {
			DummyDB = append(DummyDB[:idx], DummyDB[idx+1:]...)
		}
	}
}
