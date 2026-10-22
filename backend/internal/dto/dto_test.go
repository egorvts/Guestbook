package dto

import (
	"strings"
	"testing"
	"time"

	"github.com/egorvts/Guestbook/internal/model"
)

func TestCreateMessageRequest_ValidateAndTrim(t *testing.T) {
	tests := []struct {
		name    string
		request CreateMessageRequest
		want    bool
	}{
		{
			name: "valid request",
			request: CreateMessageRequest{
				Name:    "Egor V.",
				Content: "Hello, world!",
			},
			want: true,
		},
		{
			name: "empty name",
			request: CreateMessageRequest{
				Name:    "",
				Content: "Hello, world!",
			},
			want: false,
		},
		{
			name: "name too long",
			request: CreateMessageRequest{
				Name:    strings.Repeat("a", 256),
				Content: "Hello, world!",
			},
			want: false,
		},
		{
			name: "empty content",
			request: CreateMessageRequest{
				Name:    "Egor V.",
				Content: "",
			},
			want: false,
		},
		{
			name: "content too long",
			request: CreateMessageRequest{
				Name:    "Egor V.",
				Content: strings.Repeat("a", 2048),
			},
			want: false,
		},
		{
			name:    "whitespace-only name",
			request: CreateMessageRequest{Name: " \t\n", Content: "Hello"},
			want:    false,
		},
		{
			name:    "whitespace-only content",
			request: CreateMessageRequest{Name: "Alice", Content: " \t\n"},
			want:    false,
		},
		{
			name:    "maximum valid lengths",
			request: CreateMessageRequest{Name: strings.Repeat("a", 255), Content: strings.Repeat("b", 2047)},
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.request.ValidateAndTrim()
			if got != tt.want {
				t.Errorf("ValidateAndTrim() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromMessage(t *testing.T) {
	message := &model.Message{
		ID:        "1",
		Name:      "Egor V.",
		Content:   "Hello, world!",
		CreatedAt: time.Date(2026, time.October, 23, 0, 0, 0, 0, time.UTC),
	}

	dto := FromMessage(message)

	if dto.ID != message.ID {
		t.Errorf("FromMessage() ID = %v, want %v", dto.ID, message.ID)
	}
	if dto.Name != message.Name {
		t.Errorf("FromMessage() Name = %v, want %v", dto.Name, message.Name)
	}
	if dto.Content != message.Content {
		t.Errorf("FromMessage() Content = %v, want %v", dto.Content, message.Content)
	}
	if dto.CreatedAt != message.CreatedAt.Format(time.RFC3339) {
		t.Errorf("FromMessage() CreatedAt = %v, want %v", dto.CreatedAt, message.CreatedAt.Format(time.RFC3339))
	}
}

func TestFromMessages(t *testing.T) {
	messages := []model.Message{
		{ID: "1", Name: "Alice", Content: "Hello!", CreatedAt: time.Date(2026, time.October, 22, 12, 30, 0, 0, time.UTC)},
		{ID: "2", Name: "Bob", Content: "Hi!", CreatedAt: time.Date(2026, time.October, 22, 12, 30, 0, 0, time.UTC)},
	}

	dtos := FromMessages(messages)
	if len(dtos) != len(messages) {
		t.Fatalf("FromMessages() length = %v, want %v", len(dtos), len(messages))
	}

	for i, dto := range dtos {
		if dto.ID != messages[i].ID {
			t.Errorf("FromMessages() ID = %v, want %v", dto.ID, messages[i].ID)
		}
		if dto.Name != messages[i].Name {
			t.Errorf("FromMessages() Name = %v, want %v", dto.Name, messages[i].Name)
		}
		if dto.Content != messages[i].Content {
			t.Errorf("FromMessages() Content = %v, want %v", dto.Content, messages[i].Content)
		}
		if dto.CreatedAt != messages[i].CreatedAt.Format(time.RFC3339) {
			t.Errorf("FromMessages() CreatedAt = %v, want %v", dto.CreatedAt, messages[i].CreatedAt.Format(time.RFC3339))
		}
	}
}

func TestCreateMessageRequest_TrimsWhitespace(t *testing.T) {
	req := CreateMessageRequest{Name: " \tAlice\n", Content: "  Hello! \n"}
	if !req.ValidateAndTrim() {
		t.Fatal("expected a valid request")
	}
	if req.Name != "Alice" || req.Content != "Hello!" {
		t.Errorf("unexpected trimmed request: %+v", req)
	}
}

func TestFromMessages_Empty(t *testing.T) {
	for _, messages := range [][]model.Message{nil, {}} {
		dtos := FromMessages(messages)
		if dtos == nil || len(dtos) != 0 {
			t.Errorf("expected a non-nil empty slice, got %#v", dtos)
		}
	}
}
