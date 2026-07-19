package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/AdrianaBorges/go-legacy-integration-platform/internal/application/service"
)

type Handler struct {
	service *service.DocumentService
}

func NewHandler(service *service.DocumentService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.health)
	mux.HandleFunc("POST /api/v1/documents", h.createDocument)
	mux.HandleFunc("GET /api/v1/documents/{id}", h.getDocument)
	mux.HandleFunc("DELETE /api/v1/documents/{id}", h.deleteDocument)
	mux.HandleFunc("GET /api/v1/documents/{id}/history", h.getHistory)
	return correlationMiddleware(mux)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) createDocument(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name        string `json:"name"`
		ContentType string `json:"content_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	output, err := h.service.Create(r.Context(), service.CreateDocumentInput{
		Name:           body.Name,
		ContentType:    body.ContentType,
		IdempotencyKey: r.Header.Get("Idempotency-Key"),
	})
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	status := http.StatusCreated
	if output.Replayed {
		status = http.StatusOK
		w.Header().Set("Idempotency-Replayed", "true")
	}

	writeJSON(w, status, output.Document)
}

func (h *Handler) getDocument(w http.ResponseWriter, r *http.Request) {
	doc, err := h.service.Get(r.Context(), r.PathValue("id"))
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, doc)
}

func (h *Handler) deleteDocument(w http.ResponseWriter, r *http.Request) {
	err := h.service.Delete(r.Context(), r.PathValue("id"), correlationIDFromRequest(r))
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getHistory(w http.ResponseWriter, r *http.Request) {
	events, err := h.service.History(r.Context(), r.PathValue("id"))
	if err != nil {
		h.handleServiceError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (h *Handler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrInvalidInput):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	default:
		writeError(w, http.StatusInternalServerError, "erro interno")
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func correlationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID := strings.TrimSpace(r.Header.Get("X-Correlation-ID"))
		if correlationID == "" {
			correlationID = strings.TrimSpace(r.Header.Get("Idempotency-Key"))
		}
		if correlationID == "" {
			correlationID = "generated-by-api"
		}
		r.Header.Set("X-Correlation-ID", correlationID)
		w.Header().Set("X-Correlation-ID", correlationID)
		next.ServeHTTP(w, r)
	})
}

func correlationIDFromRequest(r *http.Request) string {
	return r.Header.Get("X-Correlation-ID")
}
