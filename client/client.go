package client

import "github.com/gofiber/websocket/v2"

type Client struct {
	Username   string
	Connection *websocket.Conn
} // Add more data to this type if needed
