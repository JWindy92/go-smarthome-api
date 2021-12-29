package websockets

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			Zap.Logger.Infow(
				"New client connected",
				"Pool size", len(pool.Clients),
				"Client", client,
			)
			break
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			Zap.Logger.Infow(
				"Client disconnected",
				"Pool size", len(pool.Clients),
				"Client", client,
			)
			break
		case message := <-pool.Broadcast:
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteMessage(message.Type, []byte(message.Body)); err != nil {
					Zap.Logger.Errorf("error broadcasting websockets message: ", err)
				}
			}
		}
	}
}
