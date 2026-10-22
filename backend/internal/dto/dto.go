package dto

import (
	"time"
	"strings"

	"github.com/egorvts/Guestbook/internal/model"
)

type CreateMessageRequest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (r *CreateMessageRequest) ValidateAndTrim() bool {
	r.Name = strings.TrimSpace(r.Name)
	r.Content = strings.TrimSpace(r.Content)
	return r.Name != "" && len(r.Name) < 256 && r.Content != "" && len(r.Content) < 2048
}

type MessageDto struct {
	ID      	string    	`json:"id"`
	Name    	string 		`json:"name"`
	CreatedAt  	string 		`json:"date"`
	Content 	string 		`json:"content"`
}

func FromMessage(m *model.Message) *MessageDto {
	return &MessageDto{
		ID:      	m.ID,
		Name:    	m.Name,
		CreatedAt: 	m.CreatedAt.Format(time.RFC3339),
		Content: 	m.Content,
	}
}

func FromMessages(messages []model.Message) []*MessageDto {
	dtos := make([]*MessageDto, len(messages))
	for i, message := range messages {
		dtos[i] = FromMessage(&message)
	}
	return dtos
}
