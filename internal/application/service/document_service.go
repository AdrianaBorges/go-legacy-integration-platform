package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/seu-usuario/go-legacy-integration-platform/internal/domain/document"
)

var (
	ErrInvalidInput = errors.New("dados do documento inválidos")
	ErrNotFound     = errors.New("documento não encontrado")
)

type DocumentRepository interface {
	Create(ctx context.Context, doc document.Document, idempotencyKey string) (document.Document, bool, error)
	Get(ctx context.Context, id string) (document.Document, error)
	Delete(ctx context.Context, id string, event document.Event) error
	History(ctx context.Context, id string) ([]document.Event, error)
}

type DocumentService struct {
	repository DocumentRepository
}

type CreateDocumentInput struct {
	Name           string
	ContentType    string
	IdempotencyKey string
}

type CreateDocumentOutput struct {
	Document document.Document
	Replayed bool
}

func NewDocumentService(repository DocumentRepository) *DocumentService {
	return &DocumentService{repository: repository}
}

func (s *DocumentService) Create(ctx context.Context, input CreateDocumentInput) (CreateDocumentOutput, error) {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.ContentType) == "" {
		return CreateDocumentOutput{}, ErrInvalidInput
	}

	now := time.Now().UTC()
	doc := document.Document{
		ID:          newID(),
		Name:        input.Name,
		ContentType: input.ContentType,
		Status:      document.StatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	created, replayed, err := s.repository.Create(ctx, doc, input.IdempotencyKey)
	if err != nil {
		return CreateDocumentOutput{}, err
	}

	return CreateDocumentOutput{Document: created, Replayed: replayed}, nil
}

func (s *DocumentService) Get(ctx context.Context, id string) (document.Document, error) {
	return s.repository.Get(ctx, id)
}

func (s *DocumentService) Delete(ctx context.Context, id, correlationID string) error {
	event := document.Event{
		ID:            newID(),
		DocumentID:    id,
		Type:          "DOCUMENT_DELETED",
		OccurredAt:    time.Now().UTC(),
		CorrelationID: correlationID,
	}
	return s.repository.Delete(ctx, id, event)
}

func (s *DocumentService) History(ctx context.Context, id string) ([]document.Event, error) {
	return s.repository.History(ctx, id)
}

func newID() string {
	buffer := make([]byte, 16)
	_, _ = rand.Read(buffer)
	return hex.EncodeToString(buffer)
}
