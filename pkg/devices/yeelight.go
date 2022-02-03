package devices

const YEELIGHT_TOPIC = "yeelight/cmnd"

// type CommandWrapper struct {
// 	Ip_addr string  `json:"ip_addr"`
// 	Cmd     Command `json:"cmd"`
// }

// type YeelightState struct {
// 	Power bool `json:"power"`
// }

// type YeelightDevice struct {
// 	Id      primitive.ObjectID `mapstructure:"_id" bson:"_id,omitempty"`
// 	Name    string             `mapstructure:"name" bson:"name"`
// 	Type    string             `mapstructure:"type" bson:"type"`
// 	Topic   string             `mapstructure:"topic" bson:"topic"`
// 	Ip_Addr string             `mapstructure:"ip_addr" bson:"ip_addr"`
// 	State   YeelightState      `mapstructure:"state" bson:"state"`
// }

// func (dev YeelightDevice) getId() primitive.ObjectID {
// 	return dev.Id
// }

// func (dev YeelightDevice) getName() string {
// 	return dev.Name
// }

// func (dev YeelightDevice) save() *mongo.InsertOneResult {
// 	m := dbutils.InitMongoInstance()
// 	defer m.Close()
// 	collection := m.Client.Database(m.Database).Collection("devices")
// 	insResult, err := collection.InsertOne(m.Context, dev)
// 	if err != nil {
// 		Zap.Logger.Errorf("error inserting new device document: %s", err)
// 	}
// 	return insResult
// }

// // TODO: This may be able to be common
// func (dev YeelightDevice) update() *mongo.UpdateResult {
// 	m := dbutils.InitMongoInstance()
// 	defer m.Close()
// 	collection := m.Client.Database(m.Database).Collection("devices")
// 	Zap.Logger.Info("Updating device")
// 	updateResult, err := collection.UpdateOne(m.Context, bson.M{"_id": dev.Id}, bson.M{"$set": dev})
// 	if err != nil {
// 		Zap.Logger.Errorf("error inserting new device document: %s", err)
// 	}
// 	Zap.Logger.Infof("Result: %d", updateResult.ModifiedCount)
// 	return updateResult
// }

// func (dev YeelightDevice) Command(command Command, mqtt_client mqtt.Client) Device {
// 	if command.validate() {
// 		wrapped := CommandWrapper{Ip_addr: dev.Ip_Addr, Cmd: command}
// 		json_cmd, err := json.Marshal(&wrapped)
// 		if err != nil {
// 			Zap.Logger.Errorf("error parsing yeelight command: %s", err)
// 		}
// 		//TODO: should return a success indicator
// 		mqtt_client.Publish(YEELIGHT_TOPIC, 1, false, json_cmd)
// 		dev.State.Power = command.powerStringToBool()
// 		dev.update()
// 	}
// 	return dev
// }
