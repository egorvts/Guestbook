package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/egorvts/Guestbook/internal/dto"
	"github.com/egorvts/Guestbook/internal/model"
	"github.com/google/uuid"
)

type mockRepository struct {
	messages    []model.Message
	err         error
	createCalls int
	created     model.Message
}

func (m *mockRepository) Create(ctx context.Context, entry *model.Message) error {
	m.createCalls++
	m.created = *entry
	return m.err
}

func (m *mockRepository) List(ctx context.Context) ([]model.Message, error) {
	return m.messages, m.err
}

func TestHandler_List(t *testing.T) {
	date := time.Date(2026, time.October, 23, 0, 0, 0, 0, time.UTC)

	repo := &mockRepository{messages: []model.Message{
		{ID: "1", Name: "Alice", Content: "Hello!", CreatedAt: date},
		{ID: "2", Name: "Bob", Content: "Hi!", CreatedAt: date},
	}}

	rr := httptest.NewRecorder()
	NewHandler(repo).List(rr, httptest.NewRequest(http.MethodGet, "/messages", nil))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if got := rr.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", got)
	}

	var messages []dto.MessageDto
	if err := json.NewDecoder(rr.Body).Decode(&messages); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	want := []dto.MessageDto{
		{ID: "1", Name: "Alice", Content: "Hello!", CreatedAt: "2026-10-23T00:00:00Z"},
		{ID: "2", Name: "Bob", Content: "Hi!", CreatedAt: "2026-10-23T00:00:00Z"},
	}
	if !reflect.DeepEqual(messages, want) {
		t.Errorf("messages = %+v, want %+v", messages, want)
	}
}

func TestHandler_List_Empty(t *testing.T) {
	rr := httptest.NewRecorder()
	NewHandler(&mockRepository{}).List(rr, httptest.NewRequest(http.MethodGet, "/messages", nil))
	
	if rr.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if got := strings.TrimSpace(rr.Body.String()); got != "[]" {
		t.Errorf("body = %q, want []", got)
	}
}

func TestHandler_List_RepositoryError(t *testing.T) {
	repo := &mockRepository{err: errors.New("database unavailable")}

	rr := httptest.NewRecorder()
	NewHandler(repo).List(rr, httptest.NewRequest(http.MethodGet, "/messages", nil))

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusInternalServerError)
	}
}

func TestHandler_Create(t *testing.T) {
	repo := &mockRepository{}

	req := httptest.NewRequest(http.MethodPost, "/messages", strings.NewReader(`{"name":"  Alice  ","content":"  Hello!  "}`))
	rr := httptest.NewRecorder()
	NewHandler(repo).Create(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusCreated)
	}
	if repo.createCalls != 1 {
		t.Fatalf("Create called %d times, want 1", repo.createCalls)
	}
	if repo.created.Name != "Alice" || repo.created.Content != "Hello!" {
		t.Errorf("unexpected saved message: %+v", repo.created)
	}
	if _, err := uuid.Parse(repo.created.ID); err != nil {
		t.Errorf("expected a UUID, got %q", repo.created.ID)
	}
}

func TestHandler_Create_InvalidRequest(t *testing.T) {
	tests := []struct{ name, body string }{
		{"empty fields", `{"name":"","content":""}`},
		{"whitespace fields", `{"name":"  ","content":"  "}`},
		{"malformed JSON", `{"name":`},
		{"empty body", ""},
		// Leading whitespace exceeds the body limit without violating field limits.
		{"oversized body", strings.Repeat(" ", MaxRequestBodySize) + `{"name":"Alice","content":"Hello!"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockRepository{}
			rr := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/messages", strings.NewReader(tt.body))
			NewHandler(repo).Create(rr, req)
			if rr.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", rr.Code, http.StatusBadRequest)
			}
			if repo.createCalls != 0 {
				t.Errorf("Create called %d times for invalid input, want 0", repo.createCalls)
			}
		})
	}
}

func TestHandler_Create_RepositoryError(t *testing.T) {
	repo := &mockRepository{err: errors.New("database unavailable")}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/messages", strings.NewReader(`{"name":"Alice","content":"Hello!"}`))

	NewHandler(repo).Create(rr, req)
	if rr.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", rr.Code, http.StatusInternalServerError)
	}
	if repo.createCalls != 1 {
		t.Errorf("Create called %d times, want 1", repo.createCalls)
	}
}
