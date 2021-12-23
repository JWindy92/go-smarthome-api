package dbutils

import (
	"context"
	"os"
	"strconv"
	"time"

	zap "github.com/JWindy92/go-smarthome-api/pkg/logwrapper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoInstance struct {
	Client   *mongo.Client
	Context  context.Context
	Cancel   context.CancelFunc
	Database string
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

var mongohost = getEnv("MONGOHOST", "10.0.0.228")

const mongoport = 27017
const dbname = "smarthome"

var Zap = zap.NewLogger()

func InitMongoInstance() MongoInstance {
	mongouri := "mongodb://" + mongohost + ":" + strconv.Itoa(mongoport)
	client, ctx, cancel, conn_err := Connect(mongouri)
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

func StringToObjectId(id string) primitive.ObjectID {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Zap.Logger.Errorf("error convrting string to ObjectId: %s", err)
	}
	return objId
}

func (m MongoInstance) Close() {
	// CancelFunc to cancel to context
	defer m.Cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {
		if err := m.Client.Disconnect(m.Context); err != nil {
			panic(err)
		}
	}()
}

func Connect(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
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
func (m MongoInstance) Ping() error {
	if err := m.Client.Ping(m.Context, readpref.Primary()); err != nil {
		return err
	}
	Zap.Logger.Info("connected successfully to mongodb")
	return nil
}

func (m MongoInstance) Query(col string, query bson.M) []bson.M {
	var devices []bson.M
	cursor, err := m.execute_query(col, query)
	if err != nil {
		Zap.Logger.Error("error querying the database", err)
	}
	err = cursor.All(m.Context, &devices)
	if err != nil {
		Zap.Logger.Error("error extracting data from mongodb cursor", err)
	}
	return devices
}

func (m MongoInstance) execute_query(col string, query bson.M) (result *mongo.Cursor, err error) {
	collection := m.Client.Database(m.Database).Collection(col)
	cursor, err := collection.Find(m.Context, query)
	if err != nil {
		Zap.Logger.Errorf("Error querying database: %v\n", err)
	}
	return cursor, err
}

// var Mongo = InitMongoInstance()
// var _ = Mongo.ping()
