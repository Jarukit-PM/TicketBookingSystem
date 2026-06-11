package tasks

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

// Enqueuer submits background tasks to asynq.
type Enqueuer interface {
	EnqueueEmailSend(ctx context.Context, bookingID string) error
}

// Client wraps asynq for enqueueing tasks from the API.
type Client struct {
	client *asynq.Client
}

// NewClient returns an asynq client using the same Redis URL as the app.
func NewClient(redisURL string) (*Client, error) {
	opt, err := redisOptFromURL(redisURL)
	if err != nil {
		return nil, err
	}
	return &Client{client: asynq.NewClient(opt)}, nil
}

// Close shuts down the underlying asynq client.
func (c *Client) Close() error {
	if c == nil || c.client == nil {
		return nil
	}
	return c.client.Close()
}

// EnqueueEmailSend enqueues a confirmation email task (worker stub until spec 09).
func (c *Client) EnqueueEmailSend(ctx context.Context, bookingID string) error {
	if c == nil || c.client == nil {
		return fmt.Errorf("task client is nil")
	}
	task, err := NewEmailSendTask(bookingID)
	if err != nil {
		return err
	}
	if _, err := c.client.EnqueueContext(ctx, task); err != nil {
		return fmt.Errorf("enqueue email:send: %w", err)
	}
	return nil
}

func redisOptFromURL(redisURL string) (asynq.RedisClientOpt, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return asynq.RedisClientOpt{}, err
	}
	return asynq.RedisClientOpt{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	}, nil
}
