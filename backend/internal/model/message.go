package model

import (
	"time"
)

type Message struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	CreatedAt time.Time `gorm:"column:date;autoCreateTime" json:"date"`
	Content   string    `gorm:"column:content" json:"content"`
}