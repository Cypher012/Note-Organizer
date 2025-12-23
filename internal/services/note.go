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

func GetNotesInFolder(db *gorm.DB, userId string, folderSlug string) (notes []models.NoteResponse, err error) {
	err = db.Model(&models.Note{}).Where("user_id = ? AND folder_slug = ?", userId, folderSlug).Find(&notes).Error

	if err != nil {
		return nil, err
	}
	return notes, nil
}

func GetNote(db *gorm.DB, userId, folderSlug, noteSlug string) (models.NoteResponse, error) {
	var folder models.Folder

	if err := db.
		Where("slug = ? AND user_id = ?", folderSlug, userId).
		First(&folder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NoteResponse{}, ErrFolderNotFound
		}
		return models.NoteResponse{}, err
	}

	var note models.NoteResponse
	if err := db.
		Model(&models.Note{}).
		Where(
			"slug = ? AND user_id = ? AND folder_slug = ?",
			noteSlug,
			userId,
			folder.Slug,
		).
		First(&note).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.NoteResponse{}, ErrNoteNotFound
		}
		return models.NoteResponse{}, err
	}

	return note, nil
}

func CreateNote(
	db *gorm.DB,
	body *models.CreateNoteRequest,
	userId,
	folderSlug string,
) error {

	var folder models.Folder
	if err := db.
		Where("slug = ? AND user_id = ?", folderSlug, userId).
		First(&folder).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFolderNotFound
		}
		return err
	}

	var existing models.Note
	err := db.
		Where(
			"user_id = ? AND folder_slug = ? AND name = ?",
			userId,
			folder.Slug,
			body.Name,
		).
		First(&existing).Error

	if err == nil {
		return ErrNoteExists
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	note := models.Note{
		ID:         uuid.NewString(),
		UserID:     userId,
		FolderSlug: folder.Slug,
		Name:       body.Name,
		Content:    body.Content,
		Slug:       slug.Make(body.Name),
	}

	return db.Create(&note).Error
}

func UpdateNote(
	db *gorm.DB,
	body *models.UpdateNoteRequest,
	userId,
	folderSlug,
	noteSlug string,
) error {

	var folder models.Folder
	if err := db.
		Where("slug = ? AND user_id = ?", folderSlug, userId).
		First(&folder).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFolderNotFound
		}
		return err
	}

	updated := map[string]any{
		"name": body.Name,
	}

	if body.Content != nil {
		updated["content"] = *body.Content
	}

	result := db.
		Model(&models.Note{}).
		Where(
			"slug = ? AND user_id = ? AND folder_slug = ?",
			noteSlug,
			userId,
			folder.Slug,
		).
		Updates(updated)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNoteNotFound
	}

	return nil
}

func DeleteNote(
	db *gorm.DB,
	userId,
	folderSlug,
	noteSlug string,
) error {

	var folder models.Folder
	if err := db.
		Where("slug = ? AND user_id = ?", folderSlug, userId).
		First(&folder).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrFolderNotFound
		}
		return err
	}

	result := db.
		Where(
			"slug = ? AND user_id = ? AND folder_slug = ?",
			noteSlug,
			userId,
			folder.Slug,
		).
		Delete(&models.Note{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNoteNotFound
	}

	return nil
}
