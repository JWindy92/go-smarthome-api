package dbutils

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	zap "github.com/JWindy92/go-smarthome-api/pkg/logwrapper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoInstance struct {
	client   *mongo.Client
	context  context.Context
	cancel   context.CancelFunc
	database string
}

type Device struct {
	Id    string `json:"Id"`
	Name  string `json:"Name"`
	State string `json:"State"`
}

const mongohost = "localhost"
const mongoport = 27017
const dbname = "smarthome"

var DummyDB []Device

var Zap = zap.NewLogger("dbutils.go")

func InitMongoInstance() MongoInstance {
	mongouri := "mongodb://" + mongohost + ":" + strconv.Itoa(mongoport)
	client, ctx, cancel, conn_err := connect(mongouri)
	if conn_err != nil {
		Zap.Logger.Panicf("error connecting to mongodb", conn_err)
	}
	mongo := MongoInstance{
		client,
		ctx,
		cancel,
		dbname,
	}
	return mongo
}

func (m MongoInstance) close() {
	// CancelFunc to cancel to context
	defer m.cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {
		if err := m.client.Disconnect(m.context); err != nil {
			panic(err)
		}
	}()
}

func connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	Zap.Logger.Infow(
		"Connecting to MongoDB",
		"host", mongohost,
		"port", mongoport,
	)
	// create context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(),
		10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	return client, ctx, cancel, err
}

// Healthcheck for mongo connection
func (m MongoInstance) ping() error {
	if err := m.client.Ping(m.context, readpref.Primary()); err != nil {
		return err
	}
	Zap.Logger.Info("connected successfully to mongodb")
	return nil
}

func (m MongoInstance) query(col string) (result *mongo.Cursor, err error) {
	collection := m.client.Database(m.database).Collection(col)
	cursor, err := collection.Find(m.context, bson.M{})
	if err != nil {
		Zap.Logger.Errorf("Error querying database: %v\n", err)
	}
	return cursor, err
}

var Mongo = InitMongoInstance()
var _ = Mongo.ping()

func GetAllDevices() []Device {
	Zap.Logger.Infow(
		"Fetching all devices",
	)
	m := InitMongoInstance()
	var devices []bson.M
	cursor, err := m.query("devices")
	if err != nil {
		Zap.Logger.Error("Idk just deal with it")
	}
	_ = cursor.All(m.context, &devices)
	fmt.Println(devices)
	return DummyDB
}

func GetDeviceById(id string) Device {
	Zap.Logger.Infow(
		"Fetching device by Id",
		"id", id,
	)
	var result Device
	for _, device := range DummyDB {
		if device.Id == id {
			result = device
		}
	}
	return result
}

func CreateNewDevice(reqBody []byte) Device {
	Zap.Logger.Infow(
		"Creating new device",
	)
	var device Device
	json.Unmarshal(reqBody, &device) // What is Unmarshal? What is '&' doing?

	DummyDB = append(DummyDB, device)

	return device
}

func DeleteDevice(id string) string {
	Zap.Logger.Infow(
		"Deleting device",
		"id", id,
	)
	for idx, device := range DummyDB {
		if device.Id == id {
			DummyDB = append(DummyDB[:idx], DummyDB[idx+1:]...)
			return device.Id
		}
	}

	return "-1"
}
