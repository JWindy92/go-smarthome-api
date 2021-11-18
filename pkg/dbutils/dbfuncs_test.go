package dbutils

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGetAllDevices(t *testing.T) {
	result := GetAllDevices()
	if len(result) < 1 {
		t.Error("No results from GetAllDevices")
	}
}

func TestGetDevicesOfType(t *testing.T) {
	sonoff := GetDevicesOfType("sonoff")
	yeelights := GetDevicesOfType("yeelight")
	invalid := GetDevicesOfType("fake_device")
	if len(sonoff) < 1 {
		t.Error("No sonoff devices returned")
	}
	if len(yeelights) < 1 {
		t.Error("No yeelight devices returned")
	}
	if len(invalid) > 0 {
		t.Error("Non existant device type returned results")
	}
}

func TestGetDeviceById(t *testing.T) {
	id, _ := primitive.ObjectIDFromHex("6193dc0ad4834c225103cabc")

	device := SonoffDevice{Id: id, Name: "sonoff1", Topic: "test/sonoff1", Type: "sonoff"}
	result := GetDeviceById(id)

	if result.getId() != device.getId() {
		t.Errorf("Returned device with _id %s, expected %s", result.getId(), device.getId())
	}
	if result.getName() != device.getName() {
		t.Errorf("Returned device with name %s, expected %s", result.getName(), device.getName())
	}
	// TODO: Once getter funcs for other attributes are in place, add checks here
}

func TestCreateAndDeleteDevice(t *testing.T) {
	prim := primitive.M{
		"name":  "test_device",
		"type":  "sonoff",
		"topic": "test/topic",
	}
	insResult := CreateNewDevice(prim)
	if insResult.InsertedID == nil {
		t.Errorf("failed to insert new device")
	}

	delResult := DeleteDevice(insResult.InsertedID.(primitive.ObjectID))
	if delResult.DeletedCount < 1 {
		t.Errorf("failed to delete device")
	}
}
