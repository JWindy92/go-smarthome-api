package websockets

import (
	"net/http"

	zap "github.com/JWindy92/go-smarthome-api/pkg/logwrapper"
	"github.com/gorilla/websocket"
)

var Zap = zap.NewLogger()

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// type socketData struct {
// 	Id    primitive.ObjectID `json:"Id"`
// 	State interface{}        `json:"State"`
// }

func Upgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Zap.Logger.Errorf("error upgrading to websocket connection: ", err)
		return nil, err
	}
	return conn, nil
}

// func wsReader(conn *websocket.Conn) {
// 	for {
// 		mt, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("read:", err)
// 			break
// 		}
// 		log.Printf("recv: %s", message)
// 		err = conn.WriteMessage(mt, message)
// 		if err != nil {
// 			log.Println("write:", err)
// 			break
// 		}
// 	}
// }

// func wsWriter(conn *websocket.Conn) {
// 	for {
// 		Zap.Logger.Info("Sending websocket msg")
// 		messageType, r, err := conn.NextReader()
// 		if err != nil {
// 			Zap.Logger.Errorf("error sending websocket message: ", err)
// 			return
// 		}
// 		w, err := conn.NextWriter(messageType)
// 		if err != nil {
// 			Zap.Logger.Errorf("error sending websocket message: ", err)
// 			return
// 		}
// 		if _, err := io.Copy(w, r); err != nil {
// 			Zap.Logger.Errorf("error sending websocket message: ", err)
// 			return
// 		}
// 		if err := w.Close(); err != nil {
// 			Zap.Logger.Errorf("error sending websocket message: ", err)
// 			return
// 		}
// 	}
// }

func ServeWs(pool *Pool, w http.ResponseWriter, r *http.Request) {
	// TODO: Actually check origin
	conn, err := Upgrade(w, r)
	if err != nil {
		Zap.Logger.Errorf("error upgrading to websocket connection: ", err)
	}
	// defer ws.Close()
	Zap.Logger.Infof("Client connected")
	client := &Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
	// go wsWriter(ws)
	// wsReader(ws)
}
