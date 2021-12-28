package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	"github.com/JWindy92/go-smarthome-api/pkg/devices"
	zap "github.com/JWindy92/go-smarthome-api/pkg/logwrapper"
	"github.com/JWindy92/go-smarthome-api/pkg/mqtt_utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Zap = zap.NewLogger()

func GetDeviceHandler(w http.ResponseWriter, r *http.Request) {

	// enableCors(&w)
	query := r.URL.Query()
	fmt.Println(query.Encode())

	Zap.Logger.Infow(
		"Handling Req",
		"method", "GET",
		"route", "/devices",
		"parameters", query.Encode(),
	)
	if query.Has("type") {
		device_type := query.Get("type")
		var result = devices.GetDevicesOfType(device_type)
		json.NewEncoder(w).Encode(result)
	} else if query.Has("id") {
		id := query.Get("id")
		var result = devices.GetDeviceById(dbutils.StringToObjectId(id))
		json.NewEncoder(w).Encode(result)
	} else {
		var result = devices.GetAllDevices()

		json.NewEncoder(w).Encode(result)
	}

}

func NewDeviceHandler(w http.ResponseWriter, r *http.Request) {
	Zap.Logger.Infow(
		"Handling Req",
		"method", "POST",
		"route", "/devices",
	)

	var prim primitive.M
	// reqBody, _ := ioutil.ReadAll(r.Body)
	_ = json.NewDecoder(r.Body).Decode(&prim)

	var result = devices.CreateNewDevice(prim)
	json.NewEncoder(w).Encode(result)
}

func DeleteDeviceHandler(w http.ResponseWriter, r *http.Request) {

	query := r.URL.Query()

	if !query.Has("id") {
		Zap.Logger.Errorf("No id provided")
	} else {
		id := query.Get("id")
		Zap.Logger.Infow(
			"Handling Req",
			"method", "DELETE",
			"route", "/devices/{id}", //? Can I format this value to contain the actual Id?
			"id", id,
		)
		devices.DeleteDevice(dbutils.StringToObjectId(id))
	}
}

func DeviceCommand(w http.ResponseWriter, r *http.Request) {
	// enableCors(&w)
	query := r.URL.Query()
	var command devices.Command
	err := json.NewDecoder(r.Body).Decode(&command)
	if err != nil {
		Zap.Logger.Errorf("error inserting new device document: %s", err)
	}

	if !query.Has("id") {
		Zap.Logger.Errorf("No id provided")
	} else {
		id := query.Get("id")
		Zap.Logger.Infow(
			"Device Command",
			"method", "POST",
			"id", id,
		)
		device := devices.GetDeviceById(dbutils.StringToObjectId(id))
		device.Command(command, mqtt_utils.MqttClient)
	}
	json.NewEncoder(w).Encode("200")
}
