package models

import (
	"time"
)

type Note struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Content   *string   `json:"content"`
	UserID    string    `json:"user_id" gorm:"not null;index"`
	FolderID  string    `json:"folder_id" gorm:"not null;index:idx_folder_note_slug,unique"`
	Slug      string    `json:"slug" gorm:"not null;index:idx_folder_note_slug,unique"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateNoteRequest struct {
	Name    string  `json:"name" validate:"required"`
	Content *string `json:"content"`
}

type UpdateNoteRequest struct {
	Name    string  `json:"name" validate:"required"`
	Content *string `json:"content"`
}
type NoteResponse struct {
	Name    string  `json:"name"`
	Content *string `json:"content"`
	Slug    string  `json:"slug"`
}
