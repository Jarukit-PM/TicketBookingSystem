package ws

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

// Client is a single WebSocket connection subscribed to one showtime room.
type Client struct {
	conn *websocket.Conn
	send chan []byte
	once sync.Once
}

// NewClient wraps a upgraded WebSocket connection.
func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 16),
	}
}

// Send enqueues a message for the write pump.
func (c *Client) Send(data []byte) error {
	select {
	case c.send <- data:
		return nil
	default:
		return websocket.ErrCloseSent
	}
}

// Close shuts down the client send channel once.
func (c *Client) Close() {
	c.once.Do(func() { close(c.send) })
}

// ReadPump reads control frames until disconnect; hub must unregister on return.
func (c *Client) ReadPump(onClose func()) {
	defer onClose()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}

// WritePump drains the send channel until closed.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
