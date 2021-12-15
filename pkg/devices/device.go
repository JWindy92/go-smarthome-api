package devices

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Command struct {
	// Id    string `mapstructure:"_id" bson:"_id,omitempty"`
	Power string `mapstructure:"power" bson:"power" json:"power"`
}

type Device interface {
	getId() primitive.ObjectID
	getName() string
	save() *mongo.InsertOneResult
	Command(command Command, mqtt_client mqtt.Client)
}

func DeviceFactory(data primitive.M) (Device, error) {
	if data["type"] == "sonoff" {
		device := SonoffDevice{}
		mapstructure.Decode(data, &device)
		return device, nil
	}
	if data["type"] == "yeelight" {
		device := YeelightDevice{}
		mapstructure.Decode(data, &device)
		return device, nil
	}
	return nil, fmt.Errorf("invalid device type %s", data["type"])
}

func mapDevicesFromPrimitives(data []primitive.M) []Device {
	var devices []Device
	for _, doc := range data {
		device, err := DeviceFactory(doc)
		if err != nil {
			Zap.Logger.Error("error mapping device object: %s", err)
		}
		devices = append(devices, device)
	}
	return devices
}
