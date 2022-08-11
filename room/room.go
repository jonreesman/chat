package room

import (
	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/jonreesman/chat/model"
)

type Room struct {
	Name       string
	register   chan *model.Client
	unregister chan *model.Client
	broadcast  chan string
	clients    map[*websocket.Conn]*model.Client
}

var Rooms map[string]*Room

func RoomSetup() {
	hubNames := []string{"main", "test2", "test3"}

	Rooms = make(map[string]*Room)

	for _, name := range hubNames {
		newUUID := uuid.NewString()
		log.Printf("%s", name)
		Rooms[newUUID] = &Room{
			Name: name,
			//register:   make(chan *websocket.Conn),
			//unregister: make(chan *websocket.Conn),
			register:   make(chan *model.Client),
			unregister: make(chan *model.Client),
			broadcast:  make(chan string),
			clients:    make(map[*websocket.Conn]*model.Client),
		}
	}

	for i, hub := range Rooms {
		log.Printf("%s", i)
		go hub.runRoom()
	}
}

func (r *Room) runRoom() {
	for {
		select {
		case connection := <-r.register:
			r.clients[connection.GetConnection()] = connection
			log.Println("connection registered")

		case message := <-r.broadcast:
			log.Println("message received:", message)

			// Send the message to all clients
			for connection := range r.clients {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Println("write error:", err)

					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
					delete(r.clients, connection)
				}
			}

		case connection := <-r.unregister:
			// Remove the client from the hub
			delete(r.clients, connection.GetConnection())

			log.Println("connection unregistered")
		}
	}
}

func (r *Room) Unregister(c *model.Client) {
	r.unregister <- c
}

func (r *Room) Register(c *model.Client) {
	r.register <- c
}

func (r *Room) Broadcast(message string) {
	r.broadcast <- message
}
