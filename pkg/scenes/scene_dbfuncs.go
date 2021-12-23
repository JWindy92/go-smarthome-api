package scenes

import (
	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// func mapScenesFromPrimitives(data []primitive.M) []Scene {
// 	var devices []Scene
// 	for _, doc := range data {
// 		device, err := DeviceFactory(doc)
// 		if err != nil {
// 			Zap.Logger.Error("error mapping device object: %s", err)
// 		}
// 		devices = append(devices, device)
// 	}
// 	return devices
// }
// func DeviceFactory(data primitive.M) (Device, error) {
// 	if data["type"] == "sonoff" {
// 		device := SonoffDevice{}
// 		mapstructure.Decode(data, &device)
// 		return device, nil
// 	}
// 	if data["type"] == "yeelight" {
// 		device := YeelightDevice{}
// 		mapstructure.Decode(data, &device)
// 		return device, nil
// 	}
// 	return nil, fmt.Errorf("invalid device type %s", data["type"])
// }
func GetSceneById(id primitive.ObjectID) Scene {
	Zap.Logger.Infow(
		"Fetching device by Id",
		"id", id,
	)
	m := dbutils.InitMongoInstance()
	defer m.Close()

	data := m.Query("devices", bson.M{"_id": id})

	scene := Scene{}
	mapstructure.Decode(data[0], &scene)

	return scene
}

func GetAllScenes() []Scene {
	Zap.Logger.Infow(
		"Fetching all scenes",
	)
	m := dbutils.InitMongoInstance()

	data := m.Query("scenes", bson.M{})

	var scenes []Scene
	for _, doc := range data {
		scene := Scene{}
		mapstructure.Decode(doc, &scene)
		scenes = append(scenes, scene)
	}

	m.Close()
	return scenes
}
