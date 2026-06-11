package db_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Jarukit-PM/TicketBookingSystem/api/internal/db"
)

func TestEnsureIndexes_Idempotent(t *testing.T) {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		t.Skip("MONGO_URI not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := db.ConnectMongo(ctx, uri)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer func() {
		_ = db.DisconnectMongo(context.Background(), client)
	}()

	database := db.Database(client, uri)
	for i := 0; i < 2; i++ {
		if err := db.EnsureIndexes(ctx, database); err != nil {
			t.Fatalf("run %d: EnsureIndexes: %v", i+1, err)
		}
	}
}

func TestDatabaseNameFromURI(t *testing.T) {
	tests := []struct {
		uri  string
		want string
	}{
		{uri: "mongodb://localhost:27017/tbs", want: "tbs"},
		{uri: "mongodb://localhost:27017/tbs?authSource=admin", want: "tbs"},
		{uri: "mongodb://localhost:27017", want: "tbs"},
	}

	for _, tt := range tests {
		if got := db.DatabaseNameFromURI(tt.uri); got != tt.want {
			t.Fatalf("DatabaseNameFromURI(%q) = %q, want %q", tt.uri, got, tt.want)
		}
	}
}
