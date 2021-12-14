package devices

import (
	"fmt"

	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

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
	fmt.Printf("Commanding %s device", dev.Type)
}
