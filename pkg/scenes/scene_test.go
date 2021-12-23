package scenes

import (
	"fmt"
	"testing"

	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	"github.com/JWindy92/go-smarthome-api/pkg/devices"
	"github.com/JWindy92/go-smarthome-api/pkg/mqtt_utils"
	"github.com/tj/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var MqttClient = mqtt_utils.MqttInit()

func TestSceneSave(t *testing.T) {
	m := dbutils.InitMongoInstance()
	defer m.Close()
	sceneId := primitive.NewObjectID()
	testScene := Scene{
		Id:   sceneId,
		Name: "Bedroom ON",
		States: []SceneState{
			{
				dbutils.StringToObjectId("61b7b943bc98e93f94a4bf37"),
				devices.Command{Power: "on"},
			},
			{
				dbutils.StringToObjectId("61b8c6a255fa968bf76d665f"),
				devices.Command{Power: "on"},
			},
		},
	}

	insResult := testScene.saveNewScene()
	assert.Equal(t, insResult.InsertedID, sceneId)
}

func TestSetScene(t *testing.T) {
	testScene := Scene{
		Name: "Test Scene",
		States: []SceneState{
			{
				dbutils.StringToObjectId("61b7b943bc98e93f94a4bf37"),
				devices.Command{Power: "on"},
			},
			{
				dbutils.StringToObjectId("61b8c6a255fa968bf76d665f"),
				devices.Command{Power: "on"},
			},
		},
	}
	testScene.SetScene(MqttClient)
}

func TestGetAllScenes(t *testing.T) {
	scenes := GetAllScenes()
	for _, scene := range scenes {
		fmt.Println("{}", scene.Name)
	}
	assert.Equal(t, 2, len(scenes))
}
