package model

import (
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type Client struct {
	Username    string
	DisplayName string
	Password    string
	ID          uuid.UUID
	connection  *websocket.Conn
	AvatarURL   string
} // Add more data to this type if needed

func CreateClient(username string, ID uuid.UUID, connection *websocket.Conn) *Client {
	return &Client{
		Username:   username,
		ID:         ID,
		connection: connection,
	}
}

func (c *Client) SetConnection(connection *websocket.Conn) {
	c.connection = connection
}

func (c Client) GetConnection() *websocket.Conn {
	return c.connection
}
