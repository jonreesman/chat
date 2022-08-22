package room

import (
	"encoding/json"
	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
	"github.com/jonreesman/chat/database"
	"github.com/jonreesman/chat/model"
)

type Room struct {
	Name          string
	register      chan *model.Client
	unregister    chan *model.Client
	broadcast     chan *model.Message
	clientsByConn map[*websocket.Conn]*model.Client
	ClientsByID   map[uuid.UUID]*model.Client
}

var Rooms map[string]*Room

func RoomSetup() {
	rooms := database.GetRooms()
	Rooms = make(map[string]*Room)

	if rooms == nil {
		hubNames := []string{"main", "test2", "test3"}
		for _, name := range hubNames {
			newUUID, _ := uuid.NewUUID()
			Rooms[newUUID.String()] = &Room{
				Name: name,
			}
			room := &model.Room{
				Name: name,
				ID:   newUUID,
			}
			database.SaveRoom(room)
		}
	}

	for _, room := range rooms {
		Rooms[room.ID.String()] = &Room{
			Name:          room.Name,
			register:      make(chan *model.Client),
			unregister:    make(chan *model.Client),
			broadcast:     make(chan *model.Message),
			clientsByConn: make(map[*websocket.Conn]*model.Client),
			ClientsByID:   make(map[uuid.UUID]*model.Client),
		}
		log.Printf("Created %s\n", room.Name)
	}

	for i, hub := range Rooms {
		log.Printf("%s", i)
		go hub.runRoom()
	}
}

func CreateRoom(room *model.Room) {
	Rooms[room.ID.String()] = &Room{
		Name:          room.Name,
		register:      make(chan *model.Client),
		unregister:    make(chan *model.Client),
		broadcast:     make(chan *model.Message),
		clientsByConn: make(map[*websocket.Conn]*model.Client),
		ClientsByID:   make(map[uuid.UUID]*model.Client),
	}
	go Rooms[room.ID.String()].runRoom()
}

func (r *Room) runRoom() {
	for {
		select {
		case connection := <-r.register:
			r.clientsByConn[connection.GetConnection()] = connection
			r.ClientsByID[connection.ID] = connection
			log.Println("connection registered")

		case message := <-r.broadcast:
			log.Println("message received:", message)
			// Send the message to all clients
			for connection := range r.clientsByConn {
				m, err := json.Marshal(message)
				if err != nil {
					log.Printf("runRoom(): failed to marshal message: %v", err)
				}
				if err := connection.WriteMessage(websocket.TextMessage, m); err != nil {
					log.Println("write error:", err)

					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
					delete(r.clientsByConn, connection)
				}
			}

		case connection := <-r.unregister:
			// Remove the client from the hub
			delete(r.clientsByConn, connection.GetConnection())
			delete(r.ClientsByID, connection.ID)
			log.Println("connection unregistered")
		}
	}
}

func DeleteRoom(id string) {
	Rooms[id].unregisterAll()
	delete(Rooms, id)
}

func (r *Room) unregisterAll() {
	for _, client := range r.clientsByConn {
		r.unregister <- client
	}
}

func (r *Room) Unregister(c *model.Client) {
	r.unregister <- c
}

func (r *Room) Register(c *model.Client) {
	r.register <- c
}

func (r *Room) Broadcast(message *model.Message) {
	r.broadcast <- message
}
