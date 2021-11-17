package dbutils

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//TODO: Breakout all device definitions to new devices package
type Device interface {
	getId() string
	getName() string
	save() *mongo.InsertOneResult
}

type SonoffDevice struct {
	Id    string `mapstructure:"_id" bson:"_id,omitempty"`
	Name  string `mapstructure:"name" bson:"name"`
	Type  string `mapstructure:"type" bson:"type"`
	Topic string `mapstructure:"topic" bson:"topic"`
}

type YeelightDevice struct {
	Id      string `mapstructure:"_id" bson:"_id,omitempty"`
	Name    string `mapstructure:"name" bson:"name"`
	Type    string `mapstructure:"type" bson:"type"`
	Topic   string `mapstructure:"topic" bson:"topic"`
	Ip_Addr string `mapstructure:"ip_addr" bson:"ip_addr"`
}

func (dev SonoffDevice) getId() string {
	return dev.Id
}

func (dev SonoffDevice) getName() string {
	return dev.Name
}

func (dev SonoffDevice) save() *mongo.InsertOneResult {
	m := InitMongoInstance()
	defer m.close()
	collection := m.client.Database(m.database).Collection("devices")
	insResult, err := collection.InsertOne(m.context, dev)
	if err != nil {
		Zap.Logger.Errorf("error inserting new device document: %s", err)
	}
	return insResult
}

func (dev YeelightDevice) getId() string {
	return dev.Id
}

func (dev YeelightDevice) getName() string {
	return dev.Name
}

func (dev YeelightDevice) save() *mongo.InsertOneResult {
	m := InitMongoInstance()
	defer m.close()
	collection := m.client.Database(m.database).Collection("devices")
	insResult, err := collection.InsertOne(m.context, dev)
	if err != nil {
		Zap.Logger.Errorf("error inserting new device document: %s", err)
	}
	return insResult
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

func GetAllDevices() []Device {
	Zap.Logger.Infow(
		"Fetching all devices",
	)
	m := InitMongoInstance()

	data := m.query("devices", bson.M{})

	devices := mapDevicesFromPrimitives(data)

	m.close()
	return devices
}

func GetDevicesOfType(device_type string) []Device {
	Zap.Logger.Infow(
		"Fetching devices by type",
		"type", device_type,
	)

	m := InitMongoInstance()

	data := m.query("devices", bson.M{"type": device_type})

	devices := mapDevicesFromPrimitives(data)

	m.close()
	return devices
}

func GetDeviceById(id string) Device {
	Zap.Logger.Infow(
		"Fetching device by Id",
		"id", id,
	)
	m := InitMongoInstance()
	defer m.close()

	data := m.query("devices", bson.M{"_id": id})

	device := mapDevicesFromPrimitives(data)[0]

	return device
}

func CreateNewDevice(reqBody primitive.M) Device {
	Zap.Logger.Infow(
		"Creating new device",
	)

	device, err := DeviceFactory(reqBody)
	if err != nil {
		Zap.Logger.Error("error mapping device object: %s", err)
	}
	device.save()

	return device
}

func DeleteDevice(id string) {
	Zap.Logger.Infow(
		"Deleting device",
		"id", id,
	)
	m := InitMongoInstance()
	defer m.close()
	collection := m.client.Database(m.database).Collection("devices")
	result, err := collection.DeleteOne(m.context, bson.M{"_id": id})
	if err != nil {
		Zap.Logger.Errorf("error inserting new device document: %s", err)
	}
	Zap.Logger.Infow(
		"Successfully deleted device",
		"_id", id,
		"num_affected", result.DeletedCount,
	)
}
