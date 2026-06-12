package audit

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestActorIDPtr(t *testing.T) {
	t.Parallel()

	if got := ActorIDPtr(primitive.NilObjectID); got != nil {
		t.Fatalf("ActorIDPtr(zero) = %v, want nil", got)
	}

	id := primitive.NewObjectID()
	got := ActorIDPtr(id)
	if got == nil || *got != id {
		t.Fatalf("ActorIDPtr(id) = %v, want %v", got, id)
	}
}
