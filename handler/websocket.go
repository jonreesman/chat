package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/gofiber/websocket/v2"
	"github.com/jonreesman/chat/config"
	"github.com/jonreesman/chat/database"
	"github.com/jonreesman/chat/middleware"
	"github.com/jonreesman/chat/model"
	"github.com/jonreesman/chat/room"
)

func ConnectToRoom(c *websocket.Conn) {
	// When the function returns, unregister the client and close the connection
	var (
		tokenProvidedUserID string
		tokenProvidedRoomID string
		userToken           string
	)

	userProvidedUserID := c.Query("user_id")
	userProvidedRoomID := c.Params("id")
	roomToken, err := middleware.ParseToken(c.Query("room_token"))
	if err != nil {
		log.Printf("room handshake failed")
		c.Close()
		return
	}
	if claims, ok := roomToken.Claims.(jwt.MapClaims); ok && roomToken.Valid {
		tokenProvidedUserID = claims["user_id"].(string)
		tokenProvidedRoomID = claims["room_id"].(string)
		userToken = claims["user_token"].(string)
	} else {
		log.Printf("room token invalid")
		c.Close()
		return
	}
	parsedUserToken, err := middleware.ParseToken(userToken)
	if err != nil {
		log.Printf("room handshake failed")
		c.Close()
		return
	}
	if userProvidedUserID != tokenProvidedUserID || userProvidedRoomID != tokenProvidedRoomID || !validToken(parsedUserToken, tokenProvidedUserID) {
		log.Printf("room token does not match")
		if ok, _ := strconv.ParseBool(config.GetConfig("DEBUG")); ok {
			fmt.Println(userProvidedUserID, " ", tokenProvidedUserID)
			fmt.Println(userProvidedRoomID, " ", tokenProvidedRoomID)
		}
		c.Close()
		return
	}
	client, err := GetClientByID(userProvidedUserID)
	if err != nil {
		log.Println(err)
		return
	}

	if _, ok := room.Rooms[tokenProvidedRoomID]; !ok {
		c.Close()
		return
	}
	room := room.Rooms[tokenProvidedRoomID]
	uuidRoomID, _ := uuid.Parse(tokenProvidedRoomID)

	client.SetConnection(c)

	defer func() {
		room.Unregister(client)
		c.Close()
	}()

	// Register the client
	room.Register(client)

	messages := database.GetRoomMessages(tokenProvidedRoomID)
	for _, message := range messages {
		if check, ok := room.ClientsByID[message.UserID]; !ok {
			userInDB, err := getClientByID(message.UserID)
			if err != nil {
				log.Printf("non-existant user detected, ignoring message.")
				continue
			}
			message.User = *userInDB
		} else {
			message.User = *check
		}
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
			} else {
				log.Println(err)
			}

			return // Calls the deferred function, i.e. closes the connection on error
		}

		if messageType == websocket.TextMessage {
			r := database.FindRoom(tokenProvidedRoomID)
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
