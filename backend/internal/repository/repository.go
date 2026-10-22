package repository

import (
	"context"
	"time"

	"github.com/egorvts/Guestbook/internal/model"

	"gorm.io/gorm"
)

const MaxMessagesGet = 16

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, entry *model.Message) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.db.WithContext(ctx).Create(entry).Error
}

func (r *Repository) List(ctx context.Context) ([]model.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var messages []model.Message
	// Limit number of messages to 24 
	err := r.db.WithContext(ctx).Order("date DESC").Limit(MaxMessagesGet).Find(&messages).Error
	return messages, err
}
