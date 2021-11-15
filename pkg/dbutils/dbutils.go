package dbutils

import (
	"encoding/json"

	zap "github.com/JWindy92/go-smarthome-api/pkg/logwrapper"
)

type Device struct {
	Id    string `json:"Id"`
	Name  string `json:"Name"`
	State string `json:"State"`
}

var DummyDB []Device

var Zap = zap.NewLogger("dbutils.go")

func GetAllDevices() []Device {
	Zap.Logger.Infow(
		"Fetching all devices",
	)
	return DummyDB
}

func GetDeviceById(id string) Device {
	Zap.Logger.Infow(
		"Fetching device by Id",
		"id", id,
	)
	var result Device
	for _, device := range DummyDB {
		if device.Id == id {
			result = device
		}
	}
	return result
}

func CreateNewDevice(reqBody []byte) Device {
	Zap.Logger.Infow(
		"Creating new device",
	)
	var device Device
	json.Unmarshal(reqBody, &device) // What is Unmarshal? What is '&' doing?

	DummyDB = append(DummyDB, device)

	return device
}

func DeleteDevice(id string) string {
	Zap.Logger.Infow(
		"Deleting device",
		"id", id,
	)
	for idx, device := range DummyDB {
		if device.Id == id {
			DummyDB = append(DummyDB[:idx], DummyDB[idx+1:]...)
			return device.Id

		}
	}

	return "-1"
}
