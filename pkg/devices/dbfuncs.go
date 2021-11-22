package devices

import (
	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	zap "github.com/JWindy92/go-smarthome-api/pkg/logwrapper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var Zap = zap.NewLogger()

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
