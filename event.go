package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Event data
type Event struct {
	ID      primitive.ObjectID
	Type    string    `json:"type,omitempty"`
	Contact string    `json:"contact,omitempty"`
	Date    time.Time `json:"date,omitempty"`
	Summary string    `json:"summary,omitempty"`
	Body    string    `json:"body,omitempty"`
	Impact  string    `json:"impact,omitempty"`
	Other   string    `json:"other,omitempty"`
}
