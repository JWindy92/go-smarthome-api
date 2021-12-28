package main

import (
	"log"
	"net/http"

	"github.com/JWindy92/go-smarthome-api/pkg/handlers"
	zap "github.com/JWindy92/go-smarthome-api/pkg/logwrapper"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

var Zap = zap.NewLogger()

var upgrader = websocket.Upgrader{}

type socketData struct {
	Id    string      `json:"Id"`
	State interface{} `json:"State"`
}

func wsReader(conn *websocket.Conn) {
	for {
		var data socketData
		conn.ReadJSON(&data)

		conn.WriteJSON(data)
	}
}
func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// TODO: Actually check origin
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Zap.Logger.Errorf("error upgrading to websocket connection: ", err)
	}
	defer ws.Close()
	Zap.Logger.Infof("Client connected")

	wsReader(ws)
}

func handleRequests() {
	router := mux.NewRouter().StrictSlash(true) // What is StrictSlash?

	router.HandleFunc("/ws", wsEndpoint)

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
