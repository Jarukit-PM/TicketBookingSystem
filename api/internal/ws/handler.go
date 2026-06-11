package ws

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/inventory"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandlerDeps holds dependencies for the showtime WebSocket endpoint.
type HandlerDeps struct {
	Hub       *Hub
	Inventory *inventory.Service
}

// Showtime handles GET /ws/showtimes/:id upgrades.
func Showtime(deps HandlerDeps) gin.HandlerFunc {
	return func(c *gin.Context) {
		showtimeID, err := primitive.ObjectIDFromHex(c.Param("id"))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		snapshot, err := deps.Inventory.Snapshot(c.Request.Context(), showtimeID)
		if err != nil {
			if errors.Is(err, inventory.ErrShowtimeNotFound) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}

		client := NewClient(conn)
		roomID := showtimeID.Hex()
		deps.Hub.Register(client, roomID)

		snapshotMsg, err := json.Marshal(Message{
			Type:    EventSnapshot,
			Payload: SnapshotPayload{Snapshot: snapshot},
		})
		if err == nil {
			_ = client.Send(snapshotMsg)
		}

		go client.WritePump()
		client.ReadPump(func() {
			deps.Hub.Unregister(client)
		})
	}
}
