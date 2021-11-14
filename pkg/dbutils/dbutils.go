package dbutils

import (
	"encoding/json"
)

type Device struct {
	Id    string `json:"Id"`
	Name  string `json:"Name"`
	State string `json:"State"`
}

var DummyDB []Device

func GetAllDevices() []Device {
	return DummyDB
}

func GetDeviceById(id string) Device {
	var result Device
	for _, device := range DummyDB {
		if device.Id == id {
			result = device
		}
	}
	return result
}

func CreateNewDevice(reqBody []byte) Device {
	var device Device
	json.Unmarshal(reqBody, &device) // What is Unmarshal? What is '&' doing?

	DummyDB = append(DummyDB, device)

	return device
}

func DeleteDevice(id string) string {

	for idx, device := range DummyDB {
		if device.Id == id {
			DummyDB = append(DummyDB[:idx], DummyDB[idx+1:]...)
			return device.Id

		}
	}

	return "-1"
}
