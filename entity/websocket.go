package entity

import "github.com/gorilla/websocket"

type WebSocketConnection struct {
	Conn *websocket.Conn
}

// SendMessage sends a message through the WebSocket connection.
func (ws *WebSocketConnection) SendMessage(message []byte) error {
	return ws.Conn.WriteMessage(websocket.TextMessage, message)
}
