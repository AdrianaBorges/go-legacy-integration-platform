package main

import (
	"log"
	"net/http"

	"github.com/AdrianaBorges/go-legacy-integration-platform/internal/application/service"
	"github.com/AdrianaBorges/go-legacy-integration-platform/internal/infrastructure/repository/memory"
	"github.com/AdrianaBorges/go-legacy-integration-platform/internal/interfaces/httpapi"
)

func main() {
	repository := memory.NewDocumentRepository()
	documentService := service.NewDocumentService(repository)
	handler := httpapi.NewHandler(documentService)

	server := &http.Server{
		Addr:              ":8080",
		Handler:           handler.Routes(),
		ReadHeaderTimeout: 5 * 1e9,
	}

	log.Println("API iniciada em http://localhost:8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
