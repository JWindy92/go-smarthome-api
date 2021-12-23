package scenes

import (
	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetSceneById(id primitive.ObjectID) Scene {
	Zap.Logger.Infow(
		"Fetching device by Id",
		"id", id,
	)
	m := dbutils.InitMongoInstance()
	defer m.Close()

	data := m.Query("scenes", bson.M{"_id": id})
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
