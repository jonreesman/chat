package handlers

import (
	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/jonreesman/chat/client"
	"github.com/jonreesman/chat/room"
)

var usernames = []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15"}

var i = 0

func ConnectToRoom(c *websocket.Conn) {
	// When the function returns, unregister the client and close the connection
	roomID := c.Params("id", "test")
	if _, ok := room.Rooms[roomID]; !ok {
		c.Close()
		return
	}

	newClient := &client.Client{
		Username:   usernames[i],
		Connection: c,
	}
	i++

	defer func() {
		room.Rooms[roomID].Unregister(newClient)
		c.Close()
	}()

	// Register the client
	room.Rooms[roomID].Register(newClient)

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}

			return // Calls the deferred function, i.e. closes the connection on error
		}

		if messageType == websocket.TextMessage {
			// Broadcast the received message
			room.Rooms[roomID].Broadcast(string(message))
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}
}
