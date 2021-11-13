package main

import (
	"encoding/json"
)

type Device struct {
	Id    string `json:"Id"`
	Name  string `json:"Name"`
	State string `json:"State"`
}

var DummyDB []Device

func getAllDevices() []Device {
	return DummyDB
}

func getDeviceById(id string) Device {
	var result Device
	for _, device := range DummyDB {
		if device.Id == id {
			result = device
		}
	}
	return result
}

func createNewDevice(reqBody []byte) Device {
	var device Device
	json.Unmarshal(reqBody, &device) // What is Unmarshal? What is '&' doing?

	DummyDB = append(DummyDB, device)

	return device
}

func deleteDevice(id string) string {

	for idx, device := range DummyDB {
		if device.Id == id {
			DummyDB = append(DummyDB[:idx], DummyDB[idx+1:]...)
			return device.Id

		}
	}

	return "-1"
}
