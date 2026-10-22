package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"uuid"

	"github.com/egorvts/Guestbook/internal/dto"
	"github.com/egorvts/Guestbook/internal/model"
)

const MaxRequestBodySize = 1048576 // 1MB

type MessageRepository interface {
	Create(context.Context, *model.Message) error
	List(context.Context) ([]model.Message, error)
}

type Handler struct {
	repository MessageRepository
}

func NewHandler(repository MessageRepository) *Handler {
	return &Handler{
		repository: repository,
	}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	messages, err := h.repository.List(r.Context())
	if err != nil {
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}

	dtos := dto.FromMessages(messages)
	WriteJSONResponse(w, http.StatusOK, dtos)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)

	var req dto.CreateMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || !req.ValidateAndTrim() {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	message := model.Message{
		ID: 		uuid.New().String(),
		Name:   		req.Name,
		Content: 	req.Content,
	}

	if err := h.repository.Create(r.Context(), &message); err != nil {
		http.Error(w, "Failed to save message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
