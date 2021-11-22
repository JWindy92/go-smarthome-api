package devices

import (
	"fmt"

	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	zap "github.com/JWindy92/go-smarthome-api/pkg/logwrapper"

	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Zap = zap.NewLogger()

//TODO: Breakout all device definitions to new devices package
type Device interface {
	getId() primitive.ObjectID
	getName() string
	save() *mongo.InsertOneResult
}

type SonoffDevice struct {
	Id    primitive.ObjectID `mapstructure:"_id" bson:"_id,omitempty"`
	Name  string             `mapstructure:"name" bson:"name"`
	Type  string             `mapstructure:"type" bson:"type"`
	Topic string             `mapstructure:"topic" bson:"topic"`
}

type YeelightDevice struct {
	Id      primitive.ObjectID `mapstructure:"_id" bson:"_id,omitempty"`
	Name    string             `mapstructure:"name" bson:"name"`
	Type    string             `mapstructure:"type" bson:"type"`
	Topic   string             `mapstructure:"topic" bson:"topic"`
	Ip_Addr string             `mapstructure:"ip_addr" bson:"ip_addr"`
}

func (dev SonoffDevice) getId() primitive.ObjectID {
	return dev.Id
}

func (dev SonoffDevice) getName() string {
	return dev.Name
}

func (dev SonoffDevice) save() *mongo.InsertOneResult {
	m := dbutils.InitMongoInstance()
	defer m.Close()
	collection := m.Client.Database(m.Database).Collection("devices")
	insResult, err := collection.InsertOne(m.Context, dev)
	if err != nil {
		Zap.Logger.Errorf("error inserting new device document: %s", err)
	}
	return insResult
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
	m := dbutils.InitMongoInstance()

	data := m.Query("devices", bson.M{})

	devices := mapDevicesFromPrimitives(data)

	m.Close()
	return devices
}

func GetDevicesOfType(device_type string) []Device {
	Zap.Logger.Infow(
		"Fetching devices by type",
		"type", device_type,
	)

	m := dbutils.InitMongoInstance()

	data := m.Query("devices", bson.M{"type": device_type})

	devices := mapDevicesFromPrimitives(data)

	m.Close()
	return devices
}

func GetDeviceById(id primitive.ObjectID) Device {
	Zap.Logger.Infow(
		"Fetching device by Id",
		"id", id,
	)
	m := dbutils.InitMongoInstance()
	defer m.Close()

	data := m.Query("devices", bson.M{"_id": id})

	device := mapDevicesFromPrimitives(data)[0]

	return device
}

func CreateNewDevice(reqBody primitive.M) *mongo.InsertOneResult {
	Zap.Logger.Infow(
		"Creating new device",
	)

	device, err := DeviceFactory(reqBody)
	if err != nil {
		Zap.Logger.Error("error mapping device object: %s", err)
	}
	result := device.save()

	return result
}

func DeleteDevice(id primitive.ObjectID) *mongo.DeleteResult {
	Zap.Logger.Infow(
		"Deleting device",
		"id", id,
	)
	m := dbutils.InitMongoInstance()
	defer m.Close()
	collection := m.Client.Database(m.Database).Collection("devices")

	result, err := collection.DeleteOne(m.Context, bson.M{"_id": id})
	if err != nil {
		Zap.Logger.Errorf("error inserting new device document: %s", err)
	}
	Zap.Logger.Infow(
		"Successfully deleted device",
		"_id", id,
		"num_affected", result.DeletedCount,
	)
	return result
}
