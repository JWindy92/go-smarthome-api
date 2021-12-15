package devices

import (
	"encoding/json"

	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const YEELIGHT_TOPIC = "yeelight/cmnd"

type CommandWrapper struct {
	Ip_addr string  `json:"ip_addr"`
	Cmd     Command `json:"cmd"`
}

type YeelightDevice struct {
	Id      primitive.ObjectID `mapstructure:"_id" bson:"_id,omitempty"`
	Name    string             `mapstructure:"name" bson:"name"`
	Type    string             `mapstructure:"type" bson:"type"`
	Topic   string             `mapstructure:"topic" bson:"topic"`
	Ip_Addr string             `mapstructure:"ip_addr" bson:"ip_addr"`
}

func (dev YeelightDevice) getId() primitive.ObjectID {
	return dev.Id
}

func (dev YeelightDevice) getName() string {
	return dev.Name
}

func (dev YeelightDevice) save() *mongo.InsertOneResult {
	m := dbutils.InitMongoInstance()
	defer m.Close()
	collection := m.Client.Database(m.Database).Collection("devices")
	insResult, err := collection.InsertOne(m.Context, dev)
	if err != nil {
		Zap.Logger.Errorf("error inserting new device document: %s", err)
	}
	return insResult
}

func (dev YeelightDevice) Command(command Command, mqtt_client mqtt.Client) {
	wrapped := CommandWrapper{Ip_addr: dev.Ip_Addr, Cmd: command}
	json_cmd, err := json.Marshal(&wrapped)
	if err != nil {
		Zap.Logger.Errorf("error parsing yeelight command: %s", err)
	}
	mqtt_client.Publish(YEELIGHT_TOPIC, 1, false, json_cmd)
}
