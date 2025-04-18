package entity

import "time"

type URL struct {
	ID        string    `bson:"_id,omitempty"`
	ShortID   string    `bson:"short_id"`
	Original  string    `bson:"original_url"`
	CreatedAt time.Time `bson:"created_at"`
}
