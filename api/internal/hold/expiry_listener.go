package hold

import (
	"context"
	"log"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/audit"
)

// SeatReleasedPublisher broadcasts seat_released after a hold TTL expires.
type SeatReleasedPublisher interface {
	PublishSeatReleased(ctx context.Context, showtimeID, seatID string) error
}

// RunExpiryListener subscribes to Redis key expiry events and logs booking_timeout
// audit entries plus WebSocket seat_released for expired hold keys.
func RunExpiryListener(
	ctx context.Context,
	rdb *redis.Client,
	publisher SeatReleasedPublisher,
	auditLog *audit.Logger,
) {
	if rdb == nil {
		return
	}

	if err := rdb.ConfigSet(ctx, "notify-keyspace-events", "Ex").Err(); err != nil {
		log.Printf("hold expiry: enable keyspace notifications: %v", err)
	}

	pubsub := rdb.PSubscribe(ctx, "__keyevent@*__:expired")
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
			handleExpiredHoldKey(ctx, rdb, msg.Payload, publisher, auditLog)
		}
	}
}

func handleExpiredHoldKey(
	ctx context.Context,
	rdb *redis.Client,
	key string,
	publisher SeatReleasedPublisher,
	auditLog *audit.Logger,
) {
	if !strings.HasPrefix(key, keyPrefix) {
		return
	}

	showtimeID, seatID, ok := ParseHoldKey(key)
	if !ok {
		return
	}

	userID := FindUserIDForExpiredSeat(ctx, rdb, showtimeID, seatID)
	if auditLog != nil {
		auditLog.BookingTimeout(ctx, userID, showtimeID, seatID)
	}
	if publisher != nil {
		if err := publisher.PublishSeatReleased(ctx, showtimeID, seatID); err != nil {
			log.Printf("hold expiry: publish seat_released %s/%s: %v", showtimeID, seatID, err)
		}
	}
}
