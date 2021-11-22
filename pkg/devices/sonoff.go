package devices

import (
	"github.com/JWindy92/go-smarthome-api/pkg/dbutils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SonoffDevice struct {
	Id    primitive.ObjectID `mapstructure:"_id" bson:"_id,omitempty"`
	Name  string             `mapstructure:"name" bson:"name"`
	Type  string             `mapstructure:"type" bson:"type"`
	Topic string             `mapstructure:"topic" bson:"topic"`
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
