package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	"github.com/JWindy92/go-smarthome-api/pkg/mqtt_utils"
	"github.com/JWindy92/go-smarthome-api/pkg/scenes"
)

func GetScenesHandler(w http.ResponseWriter, r *http.Request) {
	var result = scenes.GetAllScenes()

	json.NewEncoder(w).Encode(result)
}

func SetSceneHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if !query.Has("id") {
		Zap.Logger.Errorf("No id provided")
	} else {
		id := query.Get("id")
		Zap.Logger.Infow(
			"Setting Scene",
			"method", "POST",
			"id", id,
		)
		scene := scenes.GetSceneById(dbutils.StringToObjectId(id))
		scene.SetScene(mqtt_utils.MqttClient)
	}
	json.NewEncoder(w).Encode("200")
}
