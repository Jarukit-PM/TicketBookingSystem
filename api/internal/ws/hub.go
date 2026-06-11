package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const redisChannelPrefix = "ws:showtime:"

// Hub fans out showtime seat events to local WebSocket clients and via Redis pub/sub.
type Hub struct {
	redis *redis.Client

	mu      sync.RWMutex
	rooms   map[string]map[*Client]struct{}
	clients map[*Client]string

	register   chan clientRoom
	unregister chan *Client
	local      chan localMsg
}

type clientRoom struct {
	client     *Client
	showtimeID string
}

type localMsg struct {
	showtimeID string
	data       []byte
}

// NewHub returns a hub wired to Redis for multi-instance fan-out.
func NewHub(rdb *redis.Client) *Hub {
	return &Hub{
		redis:      rdb,
		rooms:      make(map[string]map[*Client]struct{}),
		clients:    make(map[*Client]string),
		register:   make(chan clientRoom),
		unregister: make(chan *Client),
		local:      make(chan localMsg, 64),
	}
}

// Run processes hub lifecycle until ctx is cancelled.
func (h *Hub) Run(ctx context.Context) {
	if h.redis != nil {
		go h.runRedisSubscriber(ctx)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case cr := <-h.register:
			h.mu.Lock()
			if h.rooms[cr.showtimeID] == nil {
				h.rooms[cr.showtimeID] = make(map[*Client]struct{})
			}
			h.rooms[cr.showtimeID][cr.client] = struct{}{}
			h.clients[cr.client] = cr.showtimeID
			h.mu.Unlock()
		case client := <-h.unregister:
			h.removeClient(client)
		case msg := <-h.local:
			h.broadcastLocal(msg.showtimeID, msg.data)
		}
	}
}

func (h *Hub) runRedisSubscriber(ctx context.Context) {
	pubsub := h.redis.PSubscribe(ctx, redisChannelPrefix+"*")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}
			showtimeID := strings.TrimPrefix(msg.Channel, redisChannelPrefix)
			if showtimeID == "" {
				continue
			}
			h.local <- localMsg{showtimeID: showtimeID, data: []byte(msg.Payload)}
		}
	}
}

// Register adds a client to a showtime room.
func (h *Hub) Register(client *Client, showtimeID string) {
	h.register <- clientRoom{client: client, showtimeID: showtimeID}
}

// Unregister removes a client from its room.
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// Publish marshals and publishes an event to Redis for all API instances.
func (h *Hub) Publish(ctx context.Context, showtimeID string, msg Message) error {
	if h.redis == nil {
		data, err := json.Marshal(msg)
		if err != nil {
			return fmt.Errorf("marshal ws message: %w", err)
		}
		h.local <- localMsg{showtimeID: showtimeID, data: data}
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal ws message: %w", err)
	}
	return h.redis.Publish(ctx, redisChannelPrefix+showtimeID, data).Err()
}

func (h *Hub) broadcastLocal(showtimeID string, data []byte) {
	h.mu.RLock()
	room := h.rooms[showtimeID]
	clients := make([]*Client, 0, len(room))
	for client := range room {
		clients = append(clients, client)
	}
	h.mu.RUnlock()

	for _, client := range clients {
		if err := client.Send(data); err != nil {
			log.Printf("ws send: %v", err)
		}
	}
}

func (h *Hub) removeClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	showtimeID, ok := h.clients[client]
	if !ok {
		return
	}
	delete(h.clients, client)
	if room, exists := h.rooms[showtimeID]; exists {
		delete(room, client)
		if len(room) == 0 {
			delete(h.rooms, showtimeID)
		}
	}
	client.Close()
}

// PublishSeatHeld broadcasts a seat_held event after a successful hold mutation.
func (h *Hub) PublishSeatHeld(ctx context.Context, showtimeID, seatID string, expiresAt time.Time) error {
	return h.Publish(ctx, showtimeID, Message{
		Type: EventSeatHeld,
		Payload: SeatHeldPayload{
			SeatID:    seatID,
			ExpiresAt: expiresAt,
		},
	})
}

// PublishSeatReleased broadcasts a seat_released event after a successful release.
func (h *Hub) PublishSeatReleased(ctx context.Context, showtimeID, seatID string) error {
	return h.Publish(ctx, showtimeID, Message{
		Type:    EventSeatReleased,
		Payload: SeatReleasedPayload{SeatID: seatID},
	})
}

// PublishSeatSold broadcasts a seat_sold event after a successful confirm.
func (h *Hub) PublishSeatSold(ctx context.Context, showtimeID, seatID string) error {
	return h.Publish(ctx, showtimeID, Message{
		Type:    EventSeatSold,
		Payload: SeatSoldPayload{SeatID: seatID},
	})
}
