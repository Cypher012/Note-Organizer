package models

import (
	"time"
)

type Note struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"not null;index:idx_user_note,unique,priority:3"`
	UserID    string `json:"user_id" gorm:"not null;index:idx_user_note,unique,priority:1"`
	FolderID  string `json:"folder_id" gorm:"not null;index:idx_user_note,unique,priority:2"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
