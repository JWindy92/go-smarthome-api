package devices

import (
	"fmt"

	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SonoffDevice struct {
	Id    primitive.ObjectID `mapstructure:"_id" bson:"_id,omitempty"`
	Name  string             `mapstructure:"name" bson:"name"`
	Type  string             `mapstructure:"type" bson:"type"`
	Topic string             `mapstructure:"topic" bson:"topic"`
	State bool               `mapstructure:"state" bson:"state"`
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

// TODO: This may be able to be common
func (dev SonoffDevice) update() *mongo.UpdateResult {
	m := dbutils.InitMongoInstance()
	defer m.Close()
	collection := m.Client.Database(m.Database).Collection("devices")
	updateResult, err := collection.UpdateOne(m.Context, bson.M{"_id": dev.getId()}, bson.M{"$set": dev})
	if err != nil {
		Zap.Logger.Errorf("error inserting new device document: %s", err)
	}
	Zap.Logger.Infof("Result: %d", updateResult.ModifiedCount)
	return updateResult
}

func (dev SonoffDevice) Command(command Command, mqtt_client mqtt.Client) {
	if command.validate() {
		if command.Power != "" {
			dev.power(command.Power, mqtt_client)
			dev.State = command.powerStringToBool()
		}
		fmt.Print(dev.State)
		dev.update()
	}
}

func (dev SonoffDevice) power(power string, mqtt_client mqtt.Client) {
	full_topic := "cmnd/" + dev.Topic + "/POWER"
	fmt.Printf("Sending command to %s\n", full_topic)
	mqtt_client.Publish(full_topic, 1, false, power)
}

func In(str string, list []string) bool {
	for _, val := range list {
		if val == str {
			return true
		}
	}
	return false
}
