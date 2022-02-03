package devices

import (
	"encoding/json"

	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mitchellh/mapstructure"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// type Command struct {
// 	// Id    string `mapstructure:"_id" bson:"_id,omitempty"`
// 	Power string `mapstructure:"power" bson:"power" json:"power"`
// }

//TODO: this needs to be cleaned up and made more flexible
// func (cmd Command) validate() bool {
// 	if cmd.Power != "on" && cmd.Power != "off" {
// 		return false
// 	}
// 	return true
// }

// func (cmd Command) powerStringToBool() bool {
// 	if cmd.Power == "on" {
// 		return true
// 	} else {
// 		return false
// 	}
// }

type DeviceCommand struct {
	DesiredState DeviceState
	Device       Device
}

type DeviceState struct {
	Power bool `json:"power"`
}

type Device struct {
	Id     primitive.ObjectID `mapstructure:"_id,omitempty" bson:"_id,omitempty" json:"_id,omitempty"`
	Name   string             `mapstructure:"name" bson:"name" json:"name"`
	Type   string             `mapstructure:"type" bson:"type" json:"type"`
	Topic  string             `mapstructure:"topic,omitempty" bson:"topic,omitempty" json:"topic,omitempty"`
	IpAddr string             `mapstructure:"ip_addr,omitempty" bson:"ip_addr,omitempty" json:"ip_addr,omitempty"`
	State  DeviceState        `mapstructure:"state" bson:"state" json:"state"`
}

// type Device interface {
// 	getId() primitive.ObjectID
// 	getName() string
// 	save() *mongo.InsertOneResult
// 	update() *mongo.UpdateResult
// 	Command(command Command, mqtt_client mqtt.Client) Device
// }

func (dev Device) getId() primitive.ObjectID {
	return dev.Id
}

func (dev Device) getName() string {
	return dev.Name
}

func (dev Device) save() *mongo.InsertOneResult {
	m := dbutils.InitMongoInstance()
	defer m.Close()
	collection := m.Client.Database(m.Database).Collection("devices")
	insResult, err := collection.InsertOne(m.Context, dev)
	if err != nil {
		Zap.Logger.Errorf("error inserting new device document: %s", err)
	}
	return insResult
}

func (dev Device) update() *mongo.UpdateResult {
	m := dbutils.InitMongoInstance()
	defer m.Close()
	collection := m.Client.Database(m.Database).Collection("devices")
	Zap.Logger.Info("Updating device")
	updateResult, err := collection.UpdateOne(m.Context, bson.M{"_id": dev.Id}, bson.M{"$set": dev})
	if err != nil {
		Zap.Logger.Errorf("error inserting new device document: %s", err)
	}
	Zap.Logger.Infof("Result: %d", updateResult.ModifiedCount)
	return updateResult
}

func (dev Device) Command(command DeviceState, mqtt_client mqtt.Client) Device {
	//TODO:
	//Determines endpoint based on "type"
	device_command := DeviceCommand{
		DesiredState: command,
		Device:       dev,
	}

	payload, err := json.Marshal(device_command)
	if err != nil {
		Zap.Logger.Errorf("something went wrong crafting payload: %s", err)
	}

	if dev.Type == "sonoff" {
		mqtt_client.Publish("sonoff/command", 1, false, payload)
	}
	//Sends mqtt or rest containing the command (desired state)
	//and the device data so that the service can do whatever
	//it needs to do
	Zap.Logger.Info(command)
	//Then does a dev.update() to update the DB
	dev.State = command
	//TODO: Do something with UpdateResult
	_ = dev.update()

	return dev
}

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

// TODO: Is there a better name for this now?
func mapDevicesFromPrimitives(data []primitive.M) []Device {
	var devices []Device
	for _, doc := range data {
		device := Device{}
		// device, err := DeviceFactory(doc)
		mapstructure.Decode(doc, &device)
		// if err != nil {
		// 	Zap.Logger.Error("error mapping device object: %s", err)
		// }
		devices = append(devices, device)
	}
	return devices
}
