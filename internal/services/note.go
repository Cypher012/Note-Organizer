package services

import (
	"errors"

	"github.com/Cypher012/OrganizeNoteAPi/internal/models"
	"github.com/google/uuid"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

var (
	ErrNotesNotFound = errors.New("notes not found")
	ErrNoteNotFound  = errors.New("note not found")
	ErrNoteExists    = errors.New("note already exists")
)

func GetNotes(db *gorm.DB, userId string) (notes []models.NoteResponse, err error) {
	err = db.Model(&models.Note{}).Where("user_id = ?", userId).Find(&notes).Error

	if err != nil {
		return nil, err
	}
	return notes, nil
}

func GetNotesInFolder(db *gorm.DB, userId string, folderId string) (notes []models.NoteResponse, err error) {
	err = db.Model(&models.Note{}).Where("user_id = ? AND folder_id = ?", userId, folderId).Find(&notes).Error

	if err != nil {
		return nil, err
	}
	return notes, nil
}

func GetNote(db *gorm.DB, userId, folderId, noteSlug string) (note models.NoteResponse, err error) {
	err = db.Model(&models.Note{}).Where("slugs = ? AND user_id = ? AND folder_id = ?", noteSlug, userId, folderId).First(&note).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NoteResponse{}, ErrNoteNotFound
		}
		return models.NoteResponse{}, err
	}
	return note, nil
}

func CreateNote(db *gorm.DB, body *models.CreateNoteRequest, userId, folderId string) error {
	var existing models.Note

	err := db.
		Where("user_id = ? AND folder_id = ? AND name = ?", userId, folderId, body.Name).
		First(&existing).Error

	if err == nil {
		return ErrNoteExists
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	note := models.Note{
		ID:       uuid.NewString(),
		UserID:   userId,
		FolderID: folderId,
		Name:     body.Name,
		Content:  body.Content,
		Slug:     slug.Make(body.Name),
	}

	return db.Create(&note).Error
}

func UpdateNote(db *gorm.DB, body *models.UpdateNoteRequest, userId, folderId, noteSlug string) error {
	_, err := GetNote(db, userId, folderId, noteSlug)

	if err != nil {
		return err
	}

	updatedNote := map[string]string{
		"name": body.Name,
	}

	if body.Content != nil {
		updatedNote["content"] = *body.Content
	}

	if err := db.Model(&models.Note{}).Where("slug = ? AND user_id = ? AND folder_id = ?", noteSlug, userId, folderId).Updates(updatedNote).Error; err != nil {
		return err
	}

	return nil
}

func DeleteNote(db *gorm.DB, userId, folderId, noteSlug string) error {
	result := db.
		Where("slug = ? AND user_id = ? AND folder_id = ?", noteSlug, userId, folderId).
		Delete(&models.Note{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNoteNotFound
	}

	return nil
}
