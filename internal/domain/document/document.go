package document

import "time"

type Status string

const (
	StatusPending Status = "PENDING"
	StatusActive  Status = "ACTIVE"
	StatusDeleted Status = "DELETED"
)

type Document struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ContentType string    `json:"content_type"`
	Status      Status    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Event struct {
	ID            string    `json:"id"`
	DocumentID    string    `json:"document_id"`
	Type          string    `json:"type"`
	OccurredAt    time.Time `json:"occurred_at"`
	CorrelationID string    `json:"correlation_id"`
}
