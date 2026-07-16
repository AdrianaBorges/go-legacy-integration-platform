package memory

import (
	"context"
	"sync"
	"time"

	"github.com/seu-usuario/go-legacy-integration-platform/internal/application/service"
	"github.com/seu-usuario/go-legacy-integration-platform/internal/domain/document"
)

type DocumentRepository struct {
	mu             sync.RWMutex
	documents      map[string]document.Document
	history        map[string][]document.Event
	idempotencyMap map[string]string
}

func NewDocumentRepository() *DocumentRepository {
	return &DocumentRepository{
		documents:      make(map[string]document.Document),
		history:        make(map[string][]document.Event),
		idempotencyMap: make(map[string]string),
	}
}

func (r *DocumentRepository) Create(_ context.Context, doc document.Document, idempotencyKey string) (document.Document, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if idempotencyKey != "" {
		if existingID, found := r.idempotencyMap[idempotencyKey]; found {
			return r.documents[existingID], true, nil
		}
	}

	r.documents[doc.ID] = doc
	r.history[doc.ID] = append(r.history[doc.ID], document.Event{
		ID:            doc.ID + "-created",
		DocumentID:    doc.ID,
		Type:          "DOCUMENT_CREATED",
		OccurredAt:    doc.CreatedAt,
		CorrelationID: idempotencyKey,
	})

	if idempotencyKey != "" {
		r.idempotencyMap[idempotencyKey] = doc.ID
	}

	return doc, false, nil
}

func (r *DocumentRepository) Get(_ context.Context, id string) (document.Document, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	doc, found := r.documents[id]
	if !found {
		return document.Document{}, service.ErrNotFound
	}
	return doc, nil
}

func (r *DocumentRepository) Delete(_ context.Context, id string, event document.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	doc, found := r.documents[id]
	if !found {
		return service.ErrNotFound
	}

	doc.Status = document.StatusDeleted
	doc.UpdatedAt = time.Now().UTC()
	r.documents[id] = doc
	r.history[id] = append(r.history[id], event)

	return nil
}

func (r *DocumentRepository) History(_ context.Context, id string) ([]document.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if _, found := r.documents[id]; !found {
		return nil, service.ErrNotFound
	}

	events := append([]document.Event(nil), r.history[id]...)
	return events, nil
}
