package main

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Note struct {
	ID        string `gorm:"type:char(40);primaryKey" json:"id,omitempty"`
	Title     string `gorm:"varchar(255);uniqueIndex:idx_note_title,LENGTH(255);not null" json:"title,omitempty"`
	Content string `gorm:"not null" json:"content,omitempty"`
	Category  string `gorm:"varchar(100);not null" json:"category,omitempty"`
	Published bool   `gorm:"default:false;not null" json:"published"`
	createdAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01'" json:"createdAt,omitempty"`
	updatedAt time.Time `gorm:"not null;default:'1970-01-01 00:00:01';ON UPDATE CURRENT_TIMESTAMP" json:"updatedAt,omitempty"`
}

func (note *Note) beforeCreate(tx *gorm.DB)(err error){
	note.ID = uuid.New().String()
	return nil
}

type CreateNoteSchema struct {
}

type UpdateNoteSchema struct {
}