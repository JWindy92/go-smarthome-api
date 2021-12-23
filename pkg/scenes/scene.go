package scenes

import (
	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	"github.com/JWindy92/go-smarthome-api/pkg/devices"
	zap "github.com/JWindy92/go-smarthome-api/pkg/logwrapper"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Zap = zap.NewLogger()

var DB_COLLECTION = "scenes"

type SceneState struct {
	Id      primitive.ObjectID `mapstructure:"device_id" bson:"device_id"`
	Command devices.Command    `mapstructure:"command" bson:"command"`
}

type Scene struct {
	Id     primitive.ObjectID `mapstructure:"_id" bson:"_id,omitempty"`
	Name   string             `mapstructure:"name" bson:"name"`
	States []SceneState       `mapstructure:"scene_states" bson:"scene_states"`
}

func (sc Scene) saveNewScene() *mongo.InsertOneResult {
	m := dbutils.InitMongoInstance()
	defer m.Close()
	collection := m.Client.Database(m.Database).Collection(DB_COLLECTION)
	insResult, err := collection.InsertOne(m.Context, sc)
	if err != nil {
		Zap.Logger.Errorf("error inserting new scene document: %s", err)
	}
	return insResult
}

func (sc Scene) SetScene(mqtt_client mqtt.Client) {
	for _, state := range sc.States {
		device := devices.GetDeviceById(state.Id)
		device.Command(state.Command, mqtt_client)
	}
}
