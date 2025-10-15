package model

import "time"

type Event struct {
	ID          string    `json:"event_id"`
	EventType   string    `json:"event_type"`
	AggregateId string    `json:"aggregate_id"`
	Payload     string    `json:"payload"`
	CreatedAt   time.Time `json:"timestamp"`
}
