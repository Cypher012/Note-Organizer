package models

import "time"

type Folder struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"not null;index:idx_user_folder,unique,priority:2"`
	UserID    string `json:"user_id" gorm:"not null;index:idx_user_folder, unique,priority:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateFolderRequest struct {
	Name string `json:"name" validate:"required"`
}
type FolderResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
