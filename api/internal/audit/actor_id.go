package audit

import "go.mongodb.org/mongo-driver/bson/primitive"

// ActorIDPtr returns nil for the zero ObjectId so omitempty skips actorId in JSON/BSON.
func ActorIDPtr(id primitive.ObjectID) *primitive.ObjectID {
	if id.IsZero() {
		return nil
	}
	return &id
}
