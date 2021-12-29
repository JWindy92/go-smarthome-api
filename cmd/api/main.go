package main

import (
	"log"
	"net/http"

	"github.com/JWindy92/go-smarthome-api/pkg/handlers"
	zap "github.com/JWindy92/go-smarthome-api/pkg/logwrapper"
	"github.com/JWindy92/go-smarthome-api/pkg/websockets"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var Zap = zap.NewLogger()

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true) // What is StrictSlash?
	pool := websockets.NewPool()
	go pool.Start()

	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websockets.ServeWs(pool, w, r)
	})

	router.HandleFunc("/devices", handlers.NewDeviceHandler).Methods("POST")
	router.HandleFunc("/devices", handlers.DeleteDeviceHandler).Methods("DELETE")
	router.HandleFunc("/devices", handlers.GetDeviceHandler)

	router.HandleFunc("/devices/command", handlers.DeviceCommand).Methods("POST")

	router.HandleFunc("/scenes", handlers.SetSceneHandler).Methods("POST")
	router.HandleFunc("/scenes", handlers.GetScenesHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://10.0.0.228:3000"},
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

	Zap.Logger.Infof("Ready to recieve requests")
	handleRequests()
}
