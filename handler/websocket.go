package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt"
	"github.com/jonreesman/chat/config"
	"github.com/jonreesman/chat/database"
	"github.com/jonreesman/chat/model"
	"github.com/jonreesman/chat/room"
)

func parseToken(token string) (*jwt.Token, error) {
	token = strings.ReplaceAll(token, "Bearer ", "")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.GetConfig("SECRET")), nil
	})
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		fmt.Println(claims["user_id"])
		return parsedToken, nil
	} else {
		return nil, err
	}
}

func ConnectToRoom(c *websocket.Conn) {
	// When the function returns, unregister the client and close the connection
	id := c.Query("user_id")
	client, err := GetClientByID(id)
	if err != nil {
		log.Println(err)
		return
	}
	token := c.Query("Authorization")
	parsedToken, err := parseToken(token)
	if err != nil {
		log.Printf("error parsing token: %v", err)
		return
	}
	if !validToken(parsedToken, id) {
		return
	}
	roomID := c.Params("id")
	if _, ok := room.Rooms[roomID]; !ok {
		c.Close()
		return
	}
	room := room.Rooms[roomID]
	uuidRoomID, _ := uuid.Parse(roomID)

	client.SetConnection(c)

	defer func() {
		room.Unregister(client)
		c.Close()
	}()

	// Register the client
	room.Register(client)

	messages := database.GetRoomMessages(roomID)
	for _, message := range messages {
		message.User = *(room.ClientsByID[message.UserID])
		m, err := json.Marshal(message)
		if err != nil {
			log.Printf("error marshalling room history")
		}
		if err := c.WriteMessage(websocket.TextMessage, m); err != nil {
			log.Println("write error:", err)
			c.WriteMessage(websocket.CloseMessage, []byte{})
			c.Close()
		}
	}

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("read error:", err)
			}

			return // Calls the deferred function, i.e. closes the connection on error
		}

		if messageType == websocket.TextMessage {
			r := database.FindRoom(roomID)
			content := &model.Message{
				User:      *client,
				UserID:    client.ID,
				Timestamp: time.Now().Unix(),
				Content:   string(message),
				Room:      r,
				RoomID:    uuidRoomID,
			}
			database.SaveMessage(content)
			// Broadcast the received message
			room.Broadcast(content)
		} else {
			log.Println("websocket message received of type", messageType)
		}
	}
}
