package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dbutils "github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	devices "github.com/JWindy92/go-smarthome-api/pkg/devices"
	zap "github.com/JWindy92/go-smarthome-api/pkg/logwrapper"
	"github.com/JWindy92/go-smarthome-api/pkg/mqtt_utils"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var Zap = zap.NewLogger()

var MqttClient = mqtt_utils.MqttInit()

// func enableCors(w *http.ResponseWriter) {
// 	Zap.Logger.Info("Enabling CORS")
// 	(*w).Header().Set("Access-Control-Allow-Origin", "*")
// }

// TODO: the main functions called by the handlers should be run concurrently using goroutines
func getDeviceHandler(w http.ResponseWriter, r *http.Request) {

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

func newDeviceHandler(w http.ResponseWriter, r *http.Request) {
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

func deleteDeviceHandler(w http.ResponseWriter, r *http.Request) {

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

func deviceCommand(w http.ResponseWriter, r *http.Request) {
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
		device.Command(command, MqttClient)
	}
	json.NewEncoder(w).Encode("200")
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true) // What is StrictSlash?

	// headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	// credentials := handlers.AllowCredentials()
	// methods := handlers.AllowedMethods(([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}))
	// // ttl := handlers.MaxAge(3600)
	// origins := handlers.AllowedOrigins([]string{"http://localhost:3000"})

	router.HandleFunc("/devices", newDeviceHandler).Methods("POST")
	router.HandleFunc("/devices", deleteDeviceHandler).Methods("DELETE")
	router.HandleFunc("/devices", getDeviceHandler)

	router.HandleFunc("/devices/command", deviceCommand).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":5000", handler))
}

func main() {

	defer Zap.Logger.Sync() // TODO: Find out if this single call in main is sufficent

	// TODO: these values should be populated by variables
	Zap.Logger.Infow(
		"Starting API server",
		"host", "localhost",
		"port", 5000,
	)

	Zap.Logger.Infof("PLEASE WORK!!!")
	handleRequests()
}
