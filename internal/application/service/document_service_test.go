package service_test

import (
	"context"
	"testing"

	"github.com/seu-usuario/go-legacy-integration-platform/internal/application/service"
	"github.com/seu-usuario/go-legacy-integration-platform/internal/infrastructure/repository/memory"
)

func TestCreateDocumentIsIdempotent(t *testing.T) {
	repository := memory.NewDocumentRepository()
	svc := service.NewDocumentService(repository)

	first, err := svc.Create(context.Background(), service.CreateDocumentInput{
		Name:           "documento.pdf",
		ContentType:    "application/pdf",
		IdempotencyKey: "chave-001",
	})
	if err != nil {
		t.Fatalf("primeira criação falhou: %v", err)
	}

	second, err := svc.Create(context.Background(), service.CreateDocumentInput{
		Name:           "documento.pdf",
		ContentType:    "application/pdf",
		IdempotencyKey: "chave-001",
	})
	if err != nil {
		t.Fatalf("segunda criação falhou: %v", err)
	}

	if first.Document.ID != second.Document.ID {
		t.Fatalf("esperava o mesmo documento para a mesma chave de idempotência")
	}

	if !second.Replayed {
		t.Fatalf("esperava indicação de repetição idempotente")
	}
}
