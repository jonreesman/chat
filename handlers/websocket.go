package handlers

import (
	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/jonreesman/chat/room"
)

func ConnectToRoom(c *websocket.Conn) {
	// When the function returns, unregister the client and close the connection
	roomID := c.Params("id", "test")
	if _, ok := room.Rooms[roomID]; !ok {
		c.Close()
		return
	}
	defer func() {
		room.Rooms[roomID].Unregister(c)
		c.Close()
	}()

	// Register the client
	room.Rooms[roomID].Register(c)

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
