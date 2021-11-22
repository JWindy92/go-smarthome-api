package dbutils

import (
	"context"
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
	client   *mongo.Client
	context  context.Context
	cancel   context.CancelFunc
	database string
}

const mongohost = "localhost"
const mongoport = 27017
const dbname = "smarthome"

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

func StringToObjectId(id string) primitive.ObjectID {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		Zap.Logger.Errorf("error convrting string to ObjectId: %s", err)
	}
	return objId
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

func (m MongoInstance) query(col string, query bson.M) []bson.M {
	var devices []bson.M
	cursor, err := m.execute_query("devices", query)
	if err != nil {
		Zap.Logger.Error("error querying the database", err)
	}
	err = cursor.All(m.context, &devices)
	if err != nil {
		Zap.Logger.Error("error extracting data from mongodb cursor", err)
	}
	return devices
}

func (m MongoInstance) execute_query(col string, query bson.M) (result *mongo.Cursor, err error) {
	collection := m.client.Database(m.database).Collection(col)
	cursor, err := collection.Find(m.context, query)
	if err != nil {
		Zap.Logger.Errorf("Error querying database: %v\n", err)
	}
	return cursor, err
}

func Test_Me() string {
	return "tests are cool!"
}

// var Mongo = InitMongoInstance()
// var _ = Mongo.ping()
