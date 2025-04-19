package entity

import "time"

// URL represents a URL entity stored in the database.
//
// Fields:
// - ID (string): The unique identifier for the URL document in the database.
// - ShortID (string): The shortened identifier for the URL, used for redirection.
// - Original (string): The original long URL provided by the user.
// - CreatedAt (time.Time): The timestamp when the URL was created.
type URL struct {
	ID        string    `bson:"_id,omitempty"`
	ShortID   string    `bson:"short_id"`
	Original  string    `bson:"original_url"`
	CreatedAt time.Time `bson:"created_at"`
}
