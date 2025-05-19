package api

import (
	"github.com/gofiber/websocket/v2"
)

func LiveScoreWebSocketHandler(c *websocket.Conn) {
	defer c.Close()
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			break
		}
		// Echo message for now
		if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
			break
		}
	}
} 
